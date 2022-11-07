package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/open-feature/flagd/pkg/sync"
	"github.com/open-feature/open-feature-operator/apis/core/v1alpha1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
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

const (
	featureFlagConfigurationName = "featureflagconfiguration"
	featureFlagNamespaceName     = "namespace"
)

var resyncPeriod = 1 * time.Minute

type Sync struct {
	Logger       *log.Entry
	ProviderArgs sync.ProviderArgs
	client       client.Client
}

func (k *Sync) Source() string {
	return k.ProviderArgs[featureFlagConfigurationName]
}

func (k *Sync) Fetch(ctx context.Context) (string, error) {
	if k.ProviderArgs[featureFlagConfigurationName] == "" {
		k.Logger.Error("No target feature flag configuration set")
		return "{}", nil
	}

	if k.ProviderArgs[featureFlagNamespaceName] == "" {
		k.Logger.Error("No target feature flag namespace set")
		return "{}", nil
	}

	if k.client == nil {
		k.Logger.Warn("Client not initialised")
		return "{}", nil
	}

	var ff v1alpha1.FeatureFlagConfiguration
	err := k.client.Get(ctx, client.ObjectKey{
		Name:      k.ProviderArgs[featureFlagConfigurationName],
		Namespace: k.ProviderArgs[featureFlagNamespaceName],
	}, &ff)

	return ff.Spec.FeatureFlagSpec, err
}

func (k *Sync) buildConfiguration() (*rest.Config, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	var clusterConfig *rest.Config
	var err error
	if kubeconfig != "" {
		clusterConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		clusterConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	return clusterConfig, nil
}

func (k *Sync) Notify(ctx context.Context, c chan<- sync.INotify) {
	if k.ProviderArgs[featureFlagConfigurationName] == "" {
		k.Logger.Error("No target feature flag configuration set")
		return
	}
	if k.ProviderArgs[featureFlagNamespaceName] == "" {
		k.Logger.Error("No target feature flag configuration namespace set")
		return
	}
	k.Logger.Infof("Starting kubernetes sync notifier for resource %s", k.ProviderArgs["featureflagconfiguration"])

	clusterConfig, err := k.buildConfiguration()
	if err != nil {
		k.Logger.Errorf("Error building configuration: %s", err)
	}
	if err := v1alpha1.AddToScheme(scheme.Scheme); err != nil {
		k.Logger.Panic(err.Error())
	}
	k.client, err = client.New(clusterConfig, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		k.Logger.Fatalln(err)
	}
	clusterClient, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatalln(err)
	}
	resource := v1alpha1.GroupVersion.WithResource("featureflagconfigurations")
	// The created informer will not do resyncs if the given
	// defaultEventHandlerResyncPeriod is zero.
	// For more details on resync implications refer to tools/cache/shared_informer.go
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clusterClient,
		resyncPeriod, corev1.NamespaceAll, nil)
	informer := factory.ForResource(resource).Informer()

	objectKey := client.ObjectKey{
		Name:      k.ProviderArgs[featureFlagConfigurationName],
		Namespace: k.ProviderArgs[featureFlagNamespaceName],
	}
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if err := createFuncHandler(obj, objectKey, c); err != nil {
				k.Logger.Warn(err.Error())
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if err := updateFuncHandler(oldObj, newObj, objectKey, c); err != nil {
				k.Logger.Warn(err.Error())
			}
		},
		DeleteFunc: func(obj interface{}) {
			if err := deleteFuncHandler(obj, objectKey, c); err != nil {
				k.Logger.Warn(err.Error())
			}
		},
	})

	informer.Run(ctx.Done())
}

func createFuncHandler(obj interface{}, object client.ObjectKey, c chan<- sync.INotify) error {
	var ffObj v1alpha1.FeatureFlagConfiguration
	u := obj.(*unstructured.Unstructured)
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &ffObj)
	if err != nil {
		return err
	}
	if ffObj.APIVersion != fmt.Sprintf("%s/%s", v1alpha1.GroupVersion.Group, v1alpha1.GroupVersion.Version) {
		return errors.New("invalid api version")
	}
	if ffObj.Name == object.Name {
		c <- &sync.Notifier{
			Event: sync.Event[sync.DefaultEventType]{
				EventType: sync.DefaultEventTypeCreate,
			},
		}
	}
	return nil
}

func updateFuncHandler(oldObj interface{}, newObj interface{}, object client.ObjectKey, c chan<- sync.INotify) error {
	var ffOldObj v1alpha1.FeatureFlagConfiguration
	u := oldObj.(*unstructured.Unstructured)
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &ffOldObj)
	if err != nil {
		return err
	}
	if ffOldObj.APIVersion != fmt.Sprintf("%s/%s", v1alpha1.GroupVersion.Group, v1alpha1.GroupVersion.Version) {
		return errors.New("invalid api version")
	}
	var ffNewObj v1alpha1.FeatureFlagConfiguration
	u = newObj.(*unstructured.Unstructured)
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &ffNewObj)
	if err != nil {
		return err
	}
	if ffNewObj.APIVersion != fmt.Sprintf("%s/%s", v1alpha1.GroupVersion.Group, v1alpha1.GroupVersion.Version) {
		return errors.New("invalid api version")
	}
	if object.Name == ffNewObj.Name && ffOldObj.ResourceVersion != ffNewObj.ResourceVersion {
		// Only update if there is an actual featureFlagSpec change
		c <- &sync.Notifier{
			Event: sync.Event[sync.DefaultEventType]{
				EventType: sync.DefaultEventTypeModify,
			},
		}
	}
	return nil
}

func deleteFuncHandler(obj interface{}, object client.ObjectKey, c chan<- sync.INotify) error {
	var ffObj v1alpha1.FeatureFlagConfiguration
	u := obj.(*unstructured.Unstructured)
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &ffObj)
	if err != nil {
		return err
	}
	if ffObj.APIVersion != fmt.Sprintf("%s/%s", v1alpha1.GroupVersion.Group, v1alpha1.GroupVersion.Version) {
		return errors.New("invalid api version")
	}
	if ffObj.Name == object.Name {
		c <- &sync.Notifier{
			Event: sync.Event[sync.DefaultEventType]{
				EventType: sync.DefaultEventTypeDelete,
			},
		}
	}
	return nil
}
