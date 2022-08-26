package kubernetes

import (
	"context"
	"os"

	"github.com/open-feature/flagd/pkg/sync"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesSync struct {
	Logger *log.Entry
}

type FFCInterface interface {
	Projects(namespace string) FeatureFlagConfigurationInterface
}

type FFCClient struct {
	restClient rest.Interface
}

func (c *FFCClient) Projects(namespace string) FeatureFlagConfigurationInterface {
	return &FeatureFlagClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (k *KubernetesSync) Fetch(ctx context.Context) (string, error) {

	return "", nil
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
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		k.Logger.Panic(err.Error())
		os.Exit(1)
	}
	_ = informers.NewSharedInformerFactory(clientset, 0)

}
