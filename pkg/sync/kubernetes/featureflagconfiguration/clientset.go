package featureflagconfiguration

import (
	"errors"
	"reflect"
	"time"

	"github.com/open-feature/flagd/pkg/sync"
	ffv1alpha1 "github.com/open-feature/open-feature-operator/apis/core/v1alpha1"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type FFCInterface interface {
	FeatureFlagConfigurations(namespace string) Interface
}

type FFCClient struct {
	restClient rest.Interface
}

func createFuncHandler(obj interface{}, object client.ObjectKey, c chan<- sync.INotify) error {
	if reflect.TypeOf(obj) != reflect.TypeOf(&ffv1alpha1.FeatureFlagConfiguration{}) {
		return errors.New("object is not a FeatureFlagConfiguration")
	}
	if obj.(*ffv1alpha1.FeatureFlagConfiguration).Name == object.Name {
		c <- &sync.Notifier{
			Event: sync.Event[sync.DefaultEventType]{
				EventType: sync.DefaultEventTypeCreate,
			},
		}
	}
	return nil
}

func updateFuncHandler(oldObj interface{}, object client.ObjectKey, c chan<- sync.INotify) {
	if oldObj.(*ffv1alpha1.FeatureFlagConfiguration).Name == object.Name {
		c <- &sync.Notifier{
			Event: sync.Event[sync.DefaultEventType]{
				EventType: sync.DefaultEventTypeModify,
			},
		}
	}
}

func deleteFuncHandler(obj interface{}, object client.ObjectKey, c chan<- sync.INotify) {
	if obj.(*ffv1alpha1.FeatureFlagConfiguration).Name == object.Name {
		c <- &sync.Notifier{
			Event: sync.Event[sync.DefaultEventType]{
				EventType: sync.DefaultEventTypeDelete,
			},
		}
	}
}

func WatchResources(l log.Entry, clientSet FFCInterface, refreshTime time.Duration,
	object client.ObjectKey, c chan<- sync.INotify,
) {
	ns := "*"
	if object.Namespace != "" {
		ns = object.Namespace
	}
	_, ffConfigController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.FeatureFlagConfigurations(ns).List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.FeatureFlagConfigurations(ns).Watch(lo)
			},
		},
		&ffv1alpha1.FeatureFlagConfiguration{},
		refreshTime,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if err := createFuncHandler(obj, object, c); err != nil {
					l.Warn(err.Error())
				}
			},
			DeleteFunc: func(obj interface{}) {
				deleteFuncHandler(obj, object, c)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				// This indicates a change to the custom resource
				// Typically this could be anything from a status field to a spec field
				// It is important to now assertain if it is an actual flag configuration change
				updateFuncHandler(oldObj, object, c)
			},
		},
	)
	go ffConfigController.Run(wait.NeverStop)
}

func NewForConfig(config *rest.Config) (*FFCClient, error) {
	if config == nil {
		return nil, errors.New("rest config is nil")
	}
	config.ContentConfig.GroupVersion = &schema.
		GroupVersion{
		Group:   ffv1alpha1.GroupVersion.Group,
		Version: ffv1alpha1.GroupVersion.Version,
	}
	config.APIPath = "/apis"
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	client, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, err
	}

	return &FFCClient{restClient: client}, nil
}

func (c *FFCClient) FeatureFlagConfigurations(namespace string) Interface {
	return &FeatureFlagClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
