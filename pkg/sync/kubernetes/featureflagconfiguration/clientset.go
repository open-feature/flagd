package featureflagconfiguration

import (
	"github.com/open-feature/flagd/pkg/sync"
	ffv1alpha1 "github.com/open-feature/open-feature-operator/apis/core/v1alpha1"
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
	FeatureFlagConfigurations(namespace string) FeatureFlagConfigurationInterface
}

type FFCClient struct {
	restClient rest.Interface
}

func WatchResources(clientSet FFCInterface, object client.ObjectKey, c chan<- sync.INotify) {
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
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if obj.(*ffv1alpha1.FeatureFlagConfiguration).Name == object.Name {
					c <- &sync.Notifier{
						Event: sync.Event[sync.DefaultEventType]{
							sync.DefaultEventTypeCreate,
						},
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if obj.(*ffv1alpha1.FeatureFlagConfiguration).Name == object.Name {
					c <- &sync.Notifier{
						Event: sync.Event[sync.DefaultEventType]{
							sync.DefaultEventTypeDelete,
						},
					}
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				// This indicates a change to the custom resource
				// Typically this could be anything from a status field to a spec field
				// It is important to now assertain if it is an actual flag configuration change
				if oldObj.(*ffv1alpha1.FeatureFlagConfiguration).Name == object.Name {
					// Filtered on target resource
					if oldObj.(*ffv1alpha1.FeatureFlagConfiguration).Spec.FeatureFlagSpec != newObj.(*ffv1alpha1.FeatureFlagConfiguration).Spec.FeatureFlagSpec {
						// Filtered on feature flag spec
						c <- &sync.Notifier{
							Event: sync.Event[sync.DefaultEventType]{
								sync.DefaultEventTypeModify,
							},
						}
					}
				}
			},
		},
	)
	go ffConfigController.Run(wait.NeverStop)
}

func NewForConfig(c *rest.Config) (*FFCClient, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: ffv1alpha1.GroupVersion.Group, Version: ffv1alpha1.GroupVersion.Version}
	config.APIPath = "/apis"
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &FFCClient{restClient: client}, nil
}

func (c *FFCClient) FeatureFlagConfigurations(namespace string) FeatureFlagConfigurationInterface {
	return &FeatureFlagClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
