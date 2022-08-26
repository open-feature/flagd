package kubernetes

import (
	"context"

	"github.com/open-feature/open-feature-operator/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type FeatureFlagConfigurationInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.ProjectList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.Project, error)
	Create(*v1alpha1.Project) (*v1alpha1.Project, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type FeatureFlagClient struct {
	restClient rest.Interface
	ns         string
}

func (c *FeatureFlagClient) List(opts metav1.ListOptions) (*v1alpha1.ProjectList, error) {
	result := v1alpha1.ProjectList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("featureflagconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Context()).
		Into(&result)

	return &result, err
}

func (c *FeatureFlagClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.Project, error) {
	result := v1alpha1.Project{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("featureflagconfigurations").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Context()).
		Into(&result)

	return &result, err
}

func (c *FeatureFlagClient) Create(project *v1alpha1.Project) (*v1alpha1.Project, error) {
	result := v1alpha1.Project{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("featureflagconfigurations").
		Body(project).
		Do(context.Context()).
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
		Watch(context.Context())
}
