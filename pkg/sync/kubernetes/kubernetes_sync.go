package kubernetes

import (
	"context"
	"os"
	"time"

	"github.com/open-feature/flagd/pkg/sync"
	"github.com/open-feature/open-feature-operator/apis/core/v1alpha1"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesSync struct {
	Logger *log.Entry
}

type FFCInterface interface {
	FeatureFlagConfigurations(namespace string) FeatureFlagConfigurationInterface
}

type FFCClient struct {
	restClient rest.Interface
}

func (k *KubernetesSync) WatchResources(clientSet FFCInterface) cache.Store {
	projectStore, projectController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.FeatureFlagConfigurations("*").List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.FeatureFlagConfigurations("*").Watch(lo)
			},
		},
		&v1alpha1.FeatureFlagConfiguration{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)

	go projectController.Run(wait.NeverStop)
	return projectStore
}

func NewForConfig(c *rest.Config) (*FFCClient, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupVersion.Group, Version: v1alpha1.GroupVersion.Version}
	config.APIPath = "/apis"
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
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

func (k *KubernetesSync) Fetch(ctx context.Context) (string, error) {

	return "{}", nil
}

func (k *KubernetesSync) Notify(ctx context.Context, c chan<- sync.INotify) {
	k.Logger.Info("Starting kubernetes sync notifier")

	kubeconfig := os.Getenv("KUBECONFIG")

	// Create the client configuration
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		k.Logger.Panic(err.Error())
		os.Exit(1)
	}
	clientSet, err := NewForConfig(config)
	if err != nil {
		panic(err)
	}

	fstore := k.WatchResources(clientSet)

	fstore.List()

}
