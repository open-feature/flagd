package kubernetes

import (
	"context"
	"os"
	"time"

	ffv1alpha1 "github.com/open-feature/open-feature-operator/apis/core/v1alpha1"

	"github.com/open-feature/flagd/pkg/sync"
	"github.com/open-feature/flagd/pkg/sync/kubernetes/featureflagconfiguration"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesSync struct {
	Logger *log.Entry
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
	clientSet, err := featureflagconfiguration.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ffv1alpha1.AddToScheme(scheme.Scheme)

	go featureflagconfiguration.WatchResources(clientSet)

	for {
		time.Sleep(1 * time.Second)
	}
}
