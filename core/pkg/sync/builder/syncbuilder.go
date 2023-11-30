package builder

import (
	"fmt"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/core/pkg/sync/file"
	"github.com/open-feature/flagd/core/pkg/sync/kubernetes"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	regCrd  *regexp.Regexp
	regFile *regexp.Regexp
)

func init() {
	regCrd = regexp.MustCompile("^core.openfeature.dev/")
	regFile = regexp.MustCompile("^file:")
}

type ISyncBuilder interface {
	SyncFromURI(uri string, logger *logger.Logger) (sync.ISync, error)
}

type SyncBuilder struct {
	k8sClientBuilder IK8sClientBuilder
}

func NewSyncBuilder() *SyncBuilder {
	return &SyncBuilder{
		k8sClientBuilder: &KubernetesClientBuilder{},
	}
}

func (sb *SyncBuilder) SyncFromURI(uri string, logger *logger.Logger) (sync.ISync, error) {
	switch uriB := []byte(uri); {
	// filepath may be used for debugging, not recommended in deployment
	case regFile.Match(uriB):
		return file.NewFileSync(
			uri,
			logger.WithFields(
				zap.String("component", "sync"),
				zap.String("sync", "filepath"),
			),
		), nil
	case regCrd.Match(uriB):
		reader, dynamicClient, err := sb.k8sClientBuilder.GetK8sClients()
		if err != nil {
			return nil, fmt.Errorf("error creating kubernetes clients: %w", err)
		}

		return kubernetes.NewK8sSync(
			logger.WithFields(
				zap.String("component", "sync"),
				zap.String("sync", "kubernetes"),
			),
			uri,
			reader,
			dynamicClient,
		), nil
	}
	return nil, fmt.Errorf("unrecognized URI: %s", uri)
}

type IK8sClientBuilder interface {
	GetK8sClients() (client.Reader, dynamic.Interface, error)
}

type KubernetesClientBuilder struct{}

func (kcb KubernetesClientBuilder) GetK8sClients() (client.Reader, dynamic.Interface, error) {
	clusterConfig, err := k8sClusterConfig()
	if err != nil {
		return nil, nil, err
	}

	readClient, err := client.New(clusterConfig, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create readClient: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create dynamicClient: %w", err)
	}
	return readClient, dynamicClient, nil
}

// k8sClusterConfig build K8s connection config based available configurations
func k8sClusterConfig() (*rest.Config, error) {
	cfg := os.Getenv("KUBECONFIG")

	var clusterConfig *rest.Config
	var err error

	if cfg != "" {
		clusterConfig, err = clientcmd.BuildConfigFromFlags("", cfg)
		if err != nil {
			err = fmt.Errorf("error building cluster config from flags: %w", err)
		}
	} else {
		clusterConfig, err = rest.InClusterConfig()
		if err != nil {
			err = fmt.Errorf("error fetch cluster config: %w", err)
		}
	}

	if err != nil {
		return nil, err
	}

	return clusterConfig, nil
}
