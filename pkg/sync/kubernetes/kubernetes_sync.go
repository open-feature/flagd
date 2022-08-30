package kubernetes

import (
	"context"
	"os"

	"github.com/open-feature/flagd/pkg/sync"
	"github.com/open-feature/flagd/pkg/sync/kubernetes/featureflagconfiguration"
	ffv1alpha1 "github.com/open-feature/open-feature-operator/apis/core/v1alpha1"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	controllerClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Sync struct {
	Logger       *log.Entry
	ProviderArgs sync.ProviderArgs
	client       *featureflagconfiguration.FFCClient
}

func (k *Sync) Fetch(ctx context.Context) (string, error) {
	if k.ProviderArgs["featureflagconfiguration"] == "" {
		k.Logger.Info("No target feature flag configuration set")
		return "{}", nil
	}

	if k.client == nil {
		k.Logger.Warn("Client not initialised")
		return "{}", nil
	}

	config, err := k.client.FeatureFlagConfigurations(k.ProviderArgs["featureflagconfigurationnamespace"]).
		Get(k.ProviderArgs["featureflagconfiguration"], metav1.GetOptions{
			TypeMeta: metav1.TypeMeta{
				Kind:       "FeatureFlagConfiguration",
				APIVersion: "featureflag.open-feature.io/v1alpha1",
			},
		})
	if err != nil {
		return "{}", err
	}

	return config.Spec.FeatureFlagSpec, nil
}

func (k *Sync) Notify(ctx context.Context, c chan<- sync.INotify) {
	if k.ProviderArgs["featureflagconfiguration"] == "" {
		k.Logger.Info("No target feature flag configuration set")
		return
	}
	k.Logger.Infof("Starting kubernetes sync notifier for resource %s", k.ProviderArgs["featureflagconfiguration"])
	kubeconfig := os.Getenv("KUBECONFIG")

	// Create the client configuration
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		k.Logger.Panic(err.Error())
		os.Exit(1)
	}
	k.client, err = featureflagconfiguration.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	if err := ffv1alpha1.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}

	go featureflagconfiguration.WatchResources(k.client, controllerClient.ObjectKey{
		Name: k.ProviderArgs["featureflagconfiguration"],
	}, c)
}
