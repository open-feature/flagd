package builder

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	blobSync "github.com/open-feature/flagd/core/pkg/sync/blob"
	"github.com/open-feature/flagd/core/pkg/sync/file"
	"github.com/open-feature/flagd/core/pkg/sync/grpc"
	"github.com/open-feature/flagd/core/pkg/sync/grpc/credentials"
	httpSync "github.com/open-feature/flagd/core/pkg/sync/http"
	"github.com/open-feature/flagd/core/pkg/sync/kubernetes"
	"github.com/robfig/cron"
	"go.uber.org/zap"
	"gocloud.dev/blob"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	syncProviderFile       = "file"
	syncProviderFsNotify   = "fsnotify"
	syncProviderFileInfo   = "fileinfo"
	syncProviderGrpc       = "grpc"
	syncProviderKubernetes = "kubernetes"
	syncProviderHTTP       = "http"
	syncProviderGcs        = "gcs"
	syncProviderAzblob     = "azblob"
	syncProviderS3         = "s3"
)

var (
	regCrd                *regexp.Regexp
	regURL                *regexp.Regexp
	regGRPC               *regexp.Regexp
	regGRPCSecure         *regexp.Regexp
	regGRPCCustomResolver *regexp.Regexp
	regFile               *regexp.Regexp
	regGcs                *regexp.Regexp
	regAzblob             *regexp.Regexp
	regS3                 *regexp.Regexp
)

func init() {
	regCrd = regexp.MustCompile("^core.openfeature.dev/")
	regURL = regexp.MustCompile("^https?://")
	regGRPC = regexp.MustCompile("^" + grpc.Prefix)
	regGRPCSecure = regexp.MustCompile("^" + grpc.PrefixSecure)
	regGRPCCustomResolver = regexp.MustCompile("^" + grpc.SupportedScheme)
	regFile = regexp.MustCompile("^file:")
	regGcs = regexp.MustCompile("^gs://.+?/")
	regAzblob = regexp.MustCompile("^azblob://.+?/")
	regS3 = regexp.MustCompile("^s3://.+?/")
}

type ISyncBuilder interface {
	SyncFromURI(uri string, logger *logger.Logger) (sync.ISync, error)
	SyncsFromConfig(sourceConfig []sync.SourceConfig, logger *logger.Logger) ([]sync.ISync, error)
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
		return sb.newFile(uri, logger), nil
	case regCrd.Match(uriB):
		return sb.newK8s(uri, logger)
	}
	return nil, fmt.Errorf("unrecognized URI: %s", uri)
}

func (sb *SyncBuilder) SyncsFromConfig(sourceConfigs []sync.SourceConfig, logger *logger.Logger) ([]sync.ISync, error) {
	syncImpls := make([]sync.ISync, len(sourceConfigs))
	for i, syncProvider := range sourceConfigs {
		syncImpl, err := sb.syncFromConfig(syncProvider, logger)
		if err != nil {
			return nil, fmt.Errorf("could not create sync provider: %w", err)
		}
		syncImpls[i] = syncImpl
	}
	return syncImpls, nil
}

func (sb *SyncBuilder) syncFromConfig(sourceConfig sync.SourceConfig, logger *logger.Logger) (sync.ISync, error) {
	switch sourceConfig.Provider {
	case syncProviderFile:
		return sb.newFile(sourceConfig.URI, logger), nil
	case syncProviderFsNotify:
		logger.Debug(fmt.Sprintf("using fsnotify sync-provider for: %q", sourceConfig.URI))
		return sb.newFsNotify(sourceConfig.URI, logger), nil
	case syncProviderFileInfo:
		logger.Debug(fmt.Sprintf("using fileinfo sync-provider for: %q", sourceConfig.URI))
		return sb.newFileInfo(sourceConfig.URI, logger), nil
	case syncProviderKubernetes:
		logger.Debug(fmt.Sprintf("using kubernetes sync-provider for: %s", sourceConfig.URI))
		return sb.newK8s(sourceConfig.URI, logger)
	case syncProviderHTTP:
		logger.Debug(fmt.Sprintf("using remote sync-provider for: %s", sourceConfig.URI))
		return sb.newHTTP(sourceConfig, logger), nil
	case syncProviderGrpc:
		logger.Debug(fmt.Sprintf("using grpc sync-provider for: %s", sourceConfig.URI))
		return sb.newGRPC(sourceConfig, logger), nil
	case syncProviderGcs:
		logger.Debug(fmt.Sprintf("using blob sync-provider with gcs driver for: %s", sourceConfig.URI))
		return sb.newGcs(sourceConfig, logger), nil
	case syncProviderAzblob:
		logger.Debug(fmt.Sprintf("using blob sync-provider with azblob driver for: %s", sourceConfig.URI))
		return sb.newAzblob(sourceConfig, logger)
	case syncProviderS3:
		logger.Debug(fmt.Sprintf("using blob sync-provider with s3 driver for: %s", sourceConfig.URI))
		return sb.newS3(sourceConfig, logger), nil

	default:
		return nil, fmt.Errorf("invalid sync provider: %s, must be one of with "+
			"'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s' or '%s'",
			sourceConfig.Provider, syncProviderFile, syncProviderFsNotify, syncProviderFileInfo,
			syncProviderKubernetes, syncProviderHTTP, syncProviderGrpc, syncProviderGcs, syncProviderAzblob, syncProviderS3)
	}
}

// newFile returns an fsinfo sync if we are in k8s or fileinfo if not
func (sb *SyncBuilder) newFile(uri string, logger *logger.Logger) *file.Sync {
	switch os.Getenv("KUBERNETES_SERVICE_HOST") {
	case "":
		// no k8s service host env; use fileinfo
		return sb.newFileInfo(uri, logger)
	default:
		// default to fsnotify
		return sb.newFsNotify(uri, logger)
	}
}

// return a new file.Sync that uses fsnotify under the hood
func (sb *SyncBuilder) newFsNotify(uri string, logger *logger.Logger) *file.Sync {
	return file.NewFileSync(
		regFile.ReplaceAllString(uri, ""),
		file.FSNOTIFY,
		logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", syncProviderFsNotify),
		),
	)
}

// return a new file.Sync that uses os.Stat/fs.FileInfo under the hood
func (sb *SyncBuilder) newFileInfo(uri string, logger *logger.Logger) *file.Sync {
	return file.NewFileSync(
		regFile.ReplaceAllString(uri, ""),
		file.FILEINFO,
		logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", syncProviderFileInfo),
		),
	)
}

func (sb *SyncBuilder) newK8s(uri string, logger *logger.Logger) (*kubernetes.Sync, error) {
	dynamicClient, err := sb.k8sClientBuilder.GetK8sClient()
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes clients: %w", err)
	}

	return kubernetes.NewK8sSync(
		logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "kubernetes"),
		),
		regCrd.ReplaceAllString(uri, ""),
		dynamicClient,
	), nil
}

func (sb *SyncBuilder) newHTTP(config sync.SourceConfig, logger *logger.Logger) *httpSync.Sync {
	// Default to 5 seconds
	var interval uint32 = 5
	if config.Interval != 0 {
		interval = config.Interval
	}

	var client *http.Client
	if config.OAuthConfig != nil {
		oauth := clientcredentials.Config{
			ClientID:     config.OAuthConfig.ClientId,
			ClientSecret: config.OAuthConfig.ClientSecret,
			TokenURL:     config.OAuthConfig.TokenUrl,
			AuthStyle:    oauth2.AuthStyleInParams,
		}
		client = oauth.Client(context.Background())
	} else {
		client = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	return &httpSync.Sync{
		URI:    config.URI,
		Client: client,
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "remote"),
		),
		BearerToken: config.BearerToken,
		AuthHeader:  config.AuthHeader,
		Interval:    interval,
		Cron:        cron.New(),
	}
}

func (sb *SyncBuilder) newGRPC(config sync.SourceConfig, logger *logger.Logger) *grpc.Sync {
	return &grpc.Sync{
		URI: config.URI,
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "grpc"),
		),
		CredentialBuilder: &credentials.CredentialBuilder{},
		CertPath:          config.CertPath,
		ProviderID:        config.ProviderID,
		Secure:            config.TLS,
		Selector:          config.Selector,
		MaxMsgSize:        config.MaxMsgSize,
	}
}

func (sb *SyncBuilder) newGcs(config sync.SourceConfig, logger *logger.Logger) *blobSync.Sync {
	// Extract bucket uri and object name from the full URI:
	// gs://bucket/path/to/object results in gs://bucket/ as bucketUri and
	// path/to/object as an object name.
	bucketURI := regGcs.FindString(config.URI)
	objectName := regGcs.ReplaceAllString(config.URI, "")

	// Defaults to 5 seconds if interval is not set.
	var interval uint32 = 5
	if config.Interval != 0 {
		interval = config.Interval
	}

	return &blobSync.Sync{
		Bucket: bucketURI,
		Object: objectName,

		BlobURLMux: blob.DefaultURLMux(),

		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "gcs"),
		),
		Interval: interval,
		Cron:     cron.New(),
	}
}

func (sb *SyncBuilder) newAzblob(config sync.SourceConfig, logger *logger.Logger) (*blobSync.Sync, error) {
	// Required to generate the azblob service URL
	storageAccountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	if storageAccountName == "" {
		return nil, fmt.Errorf("environment variable AZURE_STORAGE_ACCOUNT not set or is blank")
	}
	if regexp.MustCompile(`\s`).MatchString(storageAccountName) {
		return nil, fmt.Errorf("environment variable AZURE_STORAGE_ACCOUNT contains whitespace")
	}

	// Extract bucket uri and object name from the full URI:
	// azblob://bucket/path/to/object results in azblob://bucket/ as bucketUri and
	// path/to/object as an object name.
	bucketURI := regAzblob.FindString(config.URI)
	objectName := regAzblob.ReplaceAllString(config.URI, "")

	// Defaults to 5 seconds if interval is not set.
	var interval uint32 = 5
	if config.Interval != 0 {
		interval = config.Interval
	}

	return &blobSync.Sync{
		Bucket: bucketURI,
		Object: objectName,

		BlobURLMux: blob.DefaultURLMux(),

		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "azblob"),
		),
		Interval: interval,
		Cron:     cron.New(),
	}, nil
}

func (sb *SyncBuilder) newS3(config sync.SourceConfig, logger *logger.Logger) *blobSync.Sync {
	// Extract bucket uri and object name from the full URI:
	// gs://bucket/path/to/object results in gs://bucket/ as bucketUri and
	// path/to/object as an object name.
	bucketURI := regS3.FindString(config.URI)
	objectName := regS3.ReplaceAllString(config.URI, "")

	// Defaults to 5 seconds if interval is not set.
	var interval uint32 = 5
	if config.Interval != 0 {
		interval = config.Interval
	}

	return &blobSync.Sync{
		Bucket: bucketURI,
		Object: objectName,

		BlobURLMux: blob.DefaultURLMux(),

		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "s3"),
		),
		Interval: interval,
		Cron:     cron.New(),
	}
}

type IK8sClientBuilder interface {
	GetK8sClient() (dynamic.Interface, error)
}

type KubernetesClientBuilder struct{}

func (kcb KubernetesClientBuilder) GetK8sClient() (dynamic.Interface, error) {
	clusterConfig, err := k8sClusterConfig()
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create dynamicClient: %w", err)
	}
	return dynamicClient, nil
}

// k8sClusterConfig build K8s connection config based available configurations
func k8sClusterConfig() (*rest.Config, error) {
	cfg := os.Getenv("KUBECONFIG")

	var clusterConfig *rest.Config
	var err error

	if cfg != "" {
		clusterConfig, err = clientcmd.BuildConfigFromFlags("", cfg)
		if err != nil {
			return nil, fmt.Errorf("error building cluster config from flags: %w", err)
		}
	} else {
		clusterConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("error fetching cluster config: %w", err)
		}
	}

	return clusterConfig, nil
}
