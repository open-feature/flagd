package kubernetes

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/open-feature/flagd/pkg/sync"
	ffv1alpha1 "github.com/open-feature/open-feature-operator/apis/core/v1alpha1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	featureFlagConfigurationName = "featureflagconfiguration"
	featureFlagNamespaceName     = "namespace"
)

var resyncPeriod time.Duration // default of 0

type Sync struct {
	Logger       *log.Entry
	ProviderArgs sync.ProviderArgs
}

func (k *Sync) Source() string {
	return k.ProviderArgs[featureFlagConfigurationName]
}

func (k *Sync) Fetch(ctx context.Context) (string, error) {
	if k.ProviderArgs[featureFlagConfigurationName] == "" {
		k.Logger.Info("No target feature flag configuration set")
		return "{}", nil
	}

	if k.ProviderArgs[featureFlagNamespaceName] == "" {
		k.Logger.Info("No target feature flag namespace set")
		return "{}", nil
	}

	// if k.client == nil {
	// 	k.Logger.Warn("Client not initialised")
	// 	return "{}", nil
	// }

	// config, err := k.client.FeatureFlagConfigurations(k.ProviderArgs[featureFlagNamespaceName]).
	// 	Get(k.ProviderArgs[featureFlagConfigurationName], metav1.GetOptions{
	// 		TypeMeta: metav1.TypeMeta{
	// 			Kind:       "FeatureFlagConfiguration",
	// 			APIVersion: "featureflag.open-feature.io/v1alpha1",
	// 		},
	// 	})
	// if err != nil {
	// 	return "{}", err
	// }

	return "{}", nil
}

func (k *Sync) Notify(ctx context.Context, c chan<- sync.INotify) {
	if k.ProviderArgs[featureFlagConfigurationName] == "" {
		k.Logger.Info("No target feature flag configuration set")
		return
	}
	if k.ProviderArgs[featureFlagNamespaceName] == "" {
		k.Logger.Info("No target feature flag configuration namespace set")
		return
	}
	k.Logger.Infof("Starting kubernetes sync notifier for resource %s", k.ProviderArgs["featureflagconfiguration"])
	kubeconfig := os.Getenv("KUBECONFIG")

	var clusterConfig *rest.Config
	var err error
	if kubeconfig != "" {
		clusterConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		clusterConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		log.Fatalln(err)
	}

	clusterClient, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatalln(err)
	}

	resource := ffv1alpha1.GroupVersion.WithResource("featureflagconfigurations")
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clusterClient, time.Minute, corev1.NamespaceAll, nil)
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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	informer.Run(ctx.Done())

}

func createFuncHandler(obj interface{}, object client.ObjectKey, c chan<- sync.INotify) error {

	var ffObj ffv1alpha1.FeatureFlagConfiguration
	u := obj.(*unstructured.Unstructured)
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &ffObj)
	if err != nil {
		return err
	}
	if reflect.TypeOf(ffObj) != reflect.TypeOf(ffv1alpha1.FeatureFlagConfiguration{}) {
		return errors.New("object is not a FeatureFlagConfiguration")
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

	var ffOldObj ffv1alpha1.FeatureFlagConfiguration
	u := oldObj.(*unstructured.Unstructured)
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &ffOldObj)
	if err != nil {
		return err
	}
	if reflect.TypeOf(ffOldObj) != reflect.TypeOf(ffv1alpha1.FeatureFlagConfiguration{}) {
		return errors.New("object is not a FeatureFlagConfiguration")
	}
	var ffNewObj ffv1alpha1.FeatureFlagConfiguration
	u = newObj.(*unstructured.Unstructured)
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &ffNewObj)
	if err != nil {
		return err
	}
	if reflect.TypeOf(newObj) != reflect.TypeOf(ffv1alpha1.FeatureFlagConfiguration{}) {
		return errors.New("new object is not a FeatureFlagConfiguration")
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
	var ffObj ffv1alpha1.FeatureFlagConfiguration
	u := obj.(*unstructured.Unstructured)
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &ffObj)
	if err != nil {
		return err
	}
	if reflect.TypeOf(ffObj) != reflect.TypeOf(ffv1alpha1.FeatureFlagConfiguration{}) {
		return errors.New("object is not a FeatureFlagConfiguration")
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
