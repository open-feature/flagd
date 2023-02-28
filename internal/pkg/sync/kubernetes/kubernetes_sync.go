package kubernetes

import (
	"context"
	"fmt"
	"os"
	"strings"
	msync "sync"
	"time"

	"github.com/open-feature/flagd/internal/pkg/logger"
	"github.com/open-feature/flagd/internal/pkg/sync"
	"github.com/open-feature/open-feature-operator/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	resyncPeriod = 1 * time.Minute
	apiVersion   = fmt.Sprintf("%s/%s", v1alpha1.GroupVersion.Group, v1alpha1.GroupVersion.Version)
)

type Sync struct {
	Logger       *logger.Logger
	ProviderArgs sync.ProviderArgs
	URI          string

	Source     string
	ready      bool
	namespace  string
	crdName    string
	readClient client.Reader
	informer   cache.SharedInformer
}

func (k *Sync) ReSync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	fetch, err := k.fetch(ctx)
	if err != nil {
		return err
	}
	dataSync <- sync.DataSync{FlagData: fetch, Source: k.Source, Type: sync.ALL}
	return nil
}

func (k *Sync) Init(ctx context.Context) error {
	var err error

	k.namespace, k.crdName, err = parseURI(k.URI)
	if err != nil {
		return err
	}

	if err := v1alpha1.AddToScheme(scheme.Scheme); err != nil {
		return err
	}
	clusterConfig, err := k8sClusterConfig()
	if err != nil {
		return err
	}

	k.readClient, err = client.New(clusterConfig, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		return err
	}

	resource := v1alpha1.GroupVersion.WithResource("featureflagconfigurations")

	// The created informer will not do resyncs if the given defaultEventHandlerResyncPeriod is zero.
	// For more details on resync implications refer to tools/cache/shared_informer.go
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicClient, resyncPeriod, k.namespace, nil)

	k.informer = factory.ForResource(resource).Informer()

	return nil
}

func (k *Sync) IsReady() bool {
	return k.ready
}

func (k *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	k.Logger.Info(fmt.Sprintf("starting kubernetes sync notifier for resource: %s", k.URI))

	// Initial fetch
	fetch, err := k.fetch(ctx)
	if err != nil {
		k.Logger.Error(fmt.Sprintf("error with the initial fetch: %s", err.Error()))
		return err
	}

	dataSync <- sync.DataSync{FlagData: fetch, Source: k.Source, Type: sync.ALL}

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
				k.Logger.Debug("new configuration created")
				msg, err := k.fetch(ctx)
				if err != nil {
					k.Logger.Error(fmt.Sprintf("error fetching after create notification: %s", err.Error()))
					continue
				}

				dataSync <- sync.DataSync{FlagData: msg, Source: k.Source, Type: sync.ALL}
			case DefaultEventTypeModify:
				k.Logger.Debug("Configuration modified")
				msg, err := k.fetch(ctx)
				if err != nil {
					k.Logger.Error(fmt.Sprintf("error fetching after update notification: %s", err.Error()))
					continue
				}

				dataSync <- sync.DataSync{FlagData: msg, Source: k.Source, Type: sync.ALL}
			case DefaultEventTypeDelete:
				k.Logger.Debug("configuration deleted")
			case DefaultEventTypeReady:
				k.Logger.Debug("notifier ready")
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
		return "", err
	}

	if exist {
		configuration, err := toFFCfg(item)
		if err != nil {
			return "", err
		}

		k.Logger.Debug(fmt.Sprintf("resource %s served from the informer cache", k.URI))
		return configuration.Spec.FeatureFlagSpec, nil
	}

	// fallback to API access - this is an informer cache miss. Could happen at the startup where cache is not filled
	var ff v1alpha1.FeatureFlagConfiguration
	err = k.readClient.Get(ctx, client.ObjectKey{
		Name:      k.crdName,
		Namespace: k.namespace,
	}, &ff)
	if err != nil {
		return "", err
	}

	k.Logger.Debug(fmt.Sprintf("resource %s served from API server", k.URI))
	return ff.Spec.FeatureFlagSpec, nil
}

func (k *Sync) notify(ctx context.Context, c chan<- INotify) {
	objectKey := client.ObjectKey{
		Name:      k.crdName,
		Namespace: k.namespace,
	}
	if _, err := k.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			k.Logger.Info(fmt.Sprintf("kube sync notifier event: add: %s %s", objectKey.Namespace, objectKey.Name))
			if err := commonHandler(obj, objectKey, DefaultEventTypeCreate, c); err != nil {
				k.Logger.Warn(err.Error())
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			k.Logger.Info(fmt.Sprintf("kube sync notifier event: update: %s %s", objectKey.Namespace, objectKey.Name))
			if err := updateFuncHandler(oldObj, newObj, objectKey, c); err != nil {
				k.Logger.Warn(err.Error())
			}
		},
		DeleteFunc: func(obj interface{}) {
			k.Logger.Info(fmt.Sprintf("kube sync notifier event: delete: %s %s", objectKey.Namespace, objectKey.Name))
			if err := commonHandler(obj, objectKey, DefaultEventTypeDelete, c); err != nil {
				k.Logger.Warn(err.Error())
			}
		},
	}); err != nil {
		k.Logger.Fatal(err.Error())
	}

	c <- &Notifier{
		Event: Event[DefaultEventType]{
			EventType: DefaultEventTypeReady,
		},
	}

	k.informer.Run(ctx.Done())
}

// commonHandler emits the desired event if and only if handler receive an object matching apiVersion and resource name
func commonHandler(obj interface{}, object client.ObjectKey, emitEvent DefaultEventType, c chan<- INotify) error {
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
func updateFuncHandler(oldObj interface{}, newObj interface{}, object client.ObjectKey, c chan<- INotify) error {
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
func toFFCfg(object interface{}) (*v1alpha1.FeatureFlagConfiguration, error) {
	u, ok := object.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("provided value is not of type *unstructured.Unstructured")
	}

	var ffObj v1alpha1.FeatureFlagConfiguration
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &ffObj)
	if err != nil {
		return nil, err
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

// k8sClusterConfig build K8s connection config based available configurations
func k8sClusterConfig() (*rest.Config, error) {
	cfg := os.Getenv("KUBECONFIG")

	var clusterConfig *rest.Config
	var err error

	if cfg != "" {
		clusterConfig, err = clientcmd.BuildConfigFromFlags("", cfg)
	} else {
		clusterConfig, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, err
	}

	return clusterConfig, nil
}
