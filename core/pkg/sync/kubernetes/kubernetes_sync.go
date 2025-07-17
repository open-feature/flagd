package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	msync "sync"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/open-feature-operator/apis/core/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
)

var (
	resyncPeriod        = 1 * time.Minute
	apiVersion          = fmt.Sprintf("%s/%s", v1beta1.GroupVersion.Group, v1beta1.GroupVersion.Version)
	featureFlagResource = v1beta1.GroupVersion.WithResource("featureflags")
)

type SyncOption func(s *Sync)

type Sync struct {
	URI string

	ready         bool
	namespace     string
	crdName       string
	logger        *logger.Logger
	dynamicClient dynamic.Interface
	informer      cache.SharedInformer
}

func NewK8sSync(
	logger *logger.Logger,
	uri string,
	dynamicClient dynamic.Interface,
) *Sync {
	return &Sync{
		logger:        logger,
		URI:           uri,
		dynamicClient: dynamicClient,
	}
}

func (k *Sync) ReSync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	fetch, err := k.fetch(ctx)
	if err != nil {
		return fmt.Errorf("unable to fetch flag configuration: %w", err)
	}
	dataSync <- sync.DataSync{FlagData: fetch, Source: k.URI}
	return nil
}

func (k *Sync) Init(_ context.Context) error {
	var err error

	k.namespace, k.crdName, err = parseURI(k.URI)
	if err != nil {
		return fmt.Errorf("unable to parse uri %s: %w", k.URI, err)
	}

	if err := v1beta1.AddToScheme(scheme.Scheme); err != nil {
		return fmt.Errorf("unable to v1beta1 types to scheme: %w", err)
	}

	// The created informer will not do resyncs if the given defaultEventHandlerResyncPeriod is zero.
	// For more details on resync implications refer to tools/cache/shared_informer.go
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(k.dynamicClient, resyncPeriod, k.namespace, nil)

	k.informer = factory.ForResource(featureFlagResource).Informer()

	return nil
}

func (k *Sync) IsReady() bool {
	return k.ready
}

func (k *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	k.logger.Info(fmt.Sprintf("starting kubernetes sync notifier for resource: %s", k.URI))

	// Initial fetch
	fetch, err := k.fetch(ctx)
	if err != nil {
		err = fmt.Errorf("error with the initial fetch: %w", err)
		k.logger.Error(err.Error())
		return err
	}

	dataSync <- sync.DataSync{FlagData: fetch, Source: k.URI}

	notifies := make(chan INotify)

	var wg msync.WaitGroup

	// Start K8s resource notifier
	wg.Add(1)
	go func() {
		defer wg.Done()
		k.notify(ctx, notifies)
	}()

	// Start notifier watcher
	wg.Add(1)
	go func() {
		defer wg.Done()
		k.watcher(ctx, notifies, dataSync)
	}()

	wg.Wait()
	return nil
}

func (k *Sync) watcher(ctx context.Context, notifies chan INotify, dataSync chan<- sync.DataSync) {
	for {
		select {
		case <-ctx.Done():
			return
		case w := <-notifies:
			switch w.GetEvent().EventType {
			case DefaultEventTypeCreate:
				k.logger.Debug("new configuration created")
				msg, err := k.fetch(ctx)
				if err != nil {
					k.logger.Error(fmt.Sprintf("error fetching after create notification: %s", err.Error()))
					continue
				}

				dataSync <- sync.DataSync{FlagData: msg, Source: k.URI}
			case DefaultEventTypeModify:
				k.logger.Debug("Configuration modified")
				msg, err := k.fetch(ctx)
				if err != nil {
					k.logger.Error(fmt.Sprintf("error fetching after update notification: %s", err.Error()))
					continue
				}

				dataSync <- sync.DataSync{FlagData: msg, Source: k.URI}
			case DefaultEventTypeDelete:
				k.logger.Debug("configuration deleted")
			case DefaultEventTypeReady:
				k.logger.Debug("notifier ready")
				k.ready = true
			}
		}
	}
}

// fetch attempts to retrieve the latest feature flag configurations
func (k *Sync) fetch(ctx context.Context) (string, error) {
	// first check the store - avoid overloading API
	item, exist, err := k.informer.GetStore().GetByKey(k.URI)
	if err != nil {
		return "", fmt.Errorf("unable to get %s from store: %w", k.URI, err)
	}

	if exist {
		configuration, err := toFFCfg(item)
		if err != nil {
			return "", err
		}

		k.logger.Debug(fmt.Sprintf("resource %s served from the informer cache", k.URI))
		return marshallFeatureFlagSpec(configuration)
	}

	// fallback to API access - this is an informer cache miss. Could happen at the startup where cache is not filled

	ffObj, err := k.dynamicClient.
		Resource(featureFlagResource).
		Namespace(k.namespace).
		Get(ctx, k.crdName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("unable to fetch FeatureFlag %s/%s: %w", k.namespace, k.crdName, err)
	}

	k.logger.Debug(fmt.Sprintf("resource %s served from API server", k.URI))

	ff, err := toFFCfg(ffObj)
	if err != nil {
		return "", fmt.Errorf("unable to convert object %s/%s to FeatureFlag: %w", k.namespace, k.crdName, err)
	}
	return marshallFeatureFlagSpec(ff)
}

func (k *Sync) notify(ctx context.Context, c chan<- INotify) {
	objectKey := types.NamespacedName{
		Name:      k.crdName,
		Namespace: k.namespace,
	}

	if _, err := k.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			k.logger.Info(fmt.Sprintf("kube sync notifier event: add: %s %s", objectKey.Namespace, objectKey.Name))
			if err := commonHandler(obj, objectKey, DefaultEventTypeCreate, c); err != nil {
				k.logger.Warn(err.Error())
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			k.logger.Info(fmt.Sprintf("kube sync notifier event: update: %s %s", objectKey.Namespace, objectKey.Name))
			if err := updateFuncHandler(oldObj, newObj, objectKey, c); err != nil {
				k.logger.Warn(err.Error())
			}
		},
		DeleteFunc: func(obj interface{}) {
			k.logger.Info(fmt.Sprintf("kube sync notifier event: delete: %s %s", objectKey.Namespace, objectKey.Name))
			if err := commonHandler(obj, objectKey, DefaultEventTypeDelete, c); err != nil {
				k.logger.Warn(err.Error())
			}
		},
	}); err != nil {
		k.logger.Fatal(err.Error())
	}

	c <- &Notifier{
		Event: Event[DefaultEventType]{
			EventType: DefaultEventTypeReady,
		},
	}

	k.informer.Run(ctx.Done())
}

// commonHandler emits the desired event if and only if handler receive an object matching apiVersion and resource name
func commonHandler(obj interface{}, object types.NamespacedName, emitEvent DefaultEventType, c chan<- INotify) error {
	ffObj, err := toFFCfg(obj)
	if err != nil {
		return err
	}

	if ffObj.APIVersion != apiVersion {
		return fmt.Errorf("invalid api version %s, expected %s", ffObj.APIVersion, apiVersion)
	}

	if ffObj.Name == object.Name {
		c <- &Notifier{
			Event: Event[DefaultEventType]{
				EventType: emitEvent,
			},
		}
	}

	return nil
}

// updateFuncHandler handles updates. Event is emitted if and only if resource name, apiVersion of old & new are equal
func updateFuncHandler(oldObj interface{}, newObj interface{}, object types.NamespacedName, c chan<- INotify) error {
	ffOldObj, err := toFFCfg(oldObj)
	if err != nil {
		return err
	}

	if ffOldObj.APIVersion != apiVersion {
		return fmt.Errorf("invalid api version %s, expected %s", ffOldObj.APIVersion, apiVersion)
	}

	ffNewObj, err := toFFCfg(newObj)
	if err != nil {
		return err
	}

	if ffNewObj.APIVersion != apiVersion {
		return fmt.Errorf("invalid api version %s, expected %s", ffNewObj.APIVersion, apiVersion)
	}

	if object.Name == ffNewObj.Name && ffOldObj.ResourceVersion != ffNewObj.ResourceVersion {
		// Only update if there is an actual featureFlagSpec change
		c <- &Notifier{
			Event: Event[DefaultEventType]{
				EventType: DefaultEventTypeModify,
			},
		}
	}
	return nil
}

// toFFCfg attempts to covert unstructured payload to configurations
func toFFCfg(object interface{}) (*v1beta1.FeatureFlag, error) {
	u, ok := object.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("provided value is not of type *unstructured.Unstructured")
	}

	var ffObj v1beta1.FeatureFlag
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &ffObj)
	if err != nil {
		return nil, fmt.Errorf("unable to convert unstructured to v1beta1.FeatureFlag: %w", err)
	}

	return &ffObj, nil
}

// parseURI parse provided uri in the format of <namespace>/<crdName> to namespace, crdName. Results in an error
// for invalid format or failed parsing
func parseURI(uri string) (string, string, error) {
	s := strings.Split(uri, "/")
	if len(s) != 2 || len(s[0]) == 0 || len(s[1]) == 0 {
		return "", "", fmt.Errorf("invalid resource uri format, expected <namespace>/<crdName> but got: %s", uri)
	}
	return s[0], s[1], nil
}

func marshallFeatureFlagSpec(ff *v1beta1.FeatureFlag) (string, error) {
	b, err := json.Marshal(ff.Spec.FlagSpec)
	if err != nil {
		return "", fmt.Errorf("failed to marshall FlagSpec: %s", err.Error())
	}
	return string(b), nil
}
