package featureflagconfiguration

import (
	"context"

	"github.com/open-feature/open-feature-operator/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type FeatureFlagConfigurationInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.FeatureFlagConfigurationList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.FeatureFlagConfiguration, error)
	Create(*v1alpha1.FeatureFlagConfiguration) (*v1alpha1.FeatureFlagConfiguration, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type FeatureFlagClient struct {
	restClient rest.Interface
	ns         string
}

func (c *FeatureFlagClient) List(opts metav1.ListOptions) (*v1alpha1.FeatureFlagConfigurationList, error) {
	result := v1alpha1.FeatureFlagConfigurationList{}
	err := c.restClient.
		Get().
		Resource("featureflagconfigurations").
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *FeatureFlagClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.FeatureFlagConfiguration, error) {
	result := v1alpha1.FeatureFlagConfiguration{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("featureflagconfigurations").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *FeatureFlagClient) Create(project *v1alpha1.FeatureFlagConfiguration) (*v1alpha1.FeatureFlagConfiguration, error) {
	result := v1alpha1.FeatureFlagConfiguration{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("featureflagconfigurations").
		Body(project).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *FeatureFlagClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("featureflagconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.Background())
}
