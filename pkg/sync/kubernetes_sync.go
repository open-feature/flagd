package sync

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesSync struct {
	Logger *log.Entry
}

func (k *KubernetesSync) Fetch(ctx context.Context) (string, error) {

	return "", nil
}

func (k *KubernetesSync) Notify(ctx context.Context, c chan<- INotify) {
	k.Logger.Info("Starting kubernetes sync notifier")

	kubeconfig := os.Getenv("KUBECONFIG")

	// Create the client configuration
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		k.Logger.Panic(err.Error())
		os.Exit(1)
	}
	_, err = kubernetes.NewForConfig(config)
	if err != nil {
		k.Logger.Panic(err.Error())
		os.Exit(1)
	}

}
