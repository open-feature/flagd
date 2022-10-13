package kubernetes

import (
	"context"
	"errors"

	"github.com/open-feature/open-feature-operator/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type Interface interface {
	List(opts metav1.ListOptions) (*v1alpha1.FeatureFlagConfigurationList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.FeatureFlagConfiguration, error)
	Create(*v1alpha1.FeatureFlagConfiguration) (*v1alpha1.FeatureFlagConfiguration, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type FeatureFlagConfigurationImpl struct {
	restClient rest.Interface
	ns         string
}

type FeatureFlagConfigurationInterface interface {
	FeatureFlagConfigurations(namespace string) Interface
}

type FeatureFlagConfigurationRestClient struct {
	restClient rest.Interface
}

func (c *FeatureFlagConfigurationImpl) List(opts metav1.ListOptions) (*v1alpha1.FeatureFlagConfigurationList, error) {
	result := v1alpha1.FeatureFlagConfigurationList{}
	err := c.restClient.
		Get().
		Resource(featureFlagConfigurationName).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *FeatureFlagConfigurationImpl) Get(name string, opts metav1.GetOptions) (*v1alpha1.FeatureFlagConfiguration, error) {
	result := v1alpha1.FeatureFlagConfiguration{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(featureFlagConfigurationName).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *FeatureFlagConfigurationImpl) Create(project *v1alpha1.FeatureFlagConfiguration) (*v1alpha1.
	FeatureFlagConfiguration, error,
) {
	result := v1alpha1.FeatureFlagConfiguration{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(featureFlagConfigurationName).
		Body(project).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *FeatureFlagConfigurationImpl) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(featureFlagConfigurationName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.Background())
}

func NewForConfig(config *rest.Config) (*FeatureFlagConfigurationRestClient, error) {
	if config == nil {
		return nil, errors.New("rest config is nil")
	}
	config.ContentConfig.GroupVersion = &schema.
		GroupVersion{
		Group:   v1alpha1.GroupVersion.Group,
		Version: v1alpha1.GroupVersion.Version,
	}
	config.APIPath = "/apis"
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	client, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, err
	}

	return &FeatureFlagConfigurationRestClient{restClient: client}, nil
}

func (c *FeatureFlagConfigurationRestClient) FeatureFlagConfigurations(namespace string) Interface {
	return &FeatureFlagConfigurationImpl{
		restClient: c.restClient,
		ns:         namespace,
	}
}
