package cmd

import (
	"os"
	"strings"
	"github.com/open-feature/flagd/pkg/runtime"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	portFlagName            = "port"
	serviceProviderFlagName = "service-provider"
	socketPathFlagName      = "socketpath"
	syncProviderFlagName    = "sync-provider"
	evaluatorFlagName       = "evaluator"
	serverCertPathFlagName  = "server-cert-path"
	serverKeyPathFlagName   = "server-key-path"
	uriFlagName             = "uri"
	bearerTokenFlagName     = "bearer-token"
)

	serverCertPath := viper.GetString(serverCertPathFlagName)
	serverKeyPath := viper.GetString(serverKeyPathFlagName)
				Port:           viper.GetInt32(portFlagName),
				ServerKeyPath:  serverKeyPath,
				ServerCertPath: serverCertPath,
			},
			GRPCService: &service.GRPCService{},
			Logger: log.WithFields(log.Fields{
				"service":   "http",
				"component": "service",
			}),
		},
		"grpc": &service.GRPCService{
			GRPCServiceConfiguration: &service.GRPCServiceConfiguration{
				Port:           viper.GetInt32(portFlagName),
				ServerKeyPath:  serverKeyPath,
				ServerCertPath: serverCertPath,
			},
			Logger: log.WithFields(log.Fields{
				"service":   "grpc",
				"component": "service",
			}),
		},
	}

	v, ok := registeredServices[name]
	if !ok {
		return nil, errors.New("no service-provider set")
	}
	log.Debugf("Using %s service-provider\n", name)
	return v, nil
}

func findSync(name string) ([]sync.ISync, error) {
	uri := viper.GetStringSlice(uriFlagName)
	bearerToken := viper.GetString(bearerTokenFlagName)
	results := make([]sync.ISync, 0, len(uri))
	for _, u := range uri {
		registeredSync := map[string]sync.ISync{
			"filepath": &sync.FilePathSync{
				URI: u,
				Logger: log.WithFields(log.Fields{
					"sync":      "filepath",
					"component": "sync",
				}),
			},
			"remote": &sync.HTTPSync{
				URI:         u,
				BearerToken: bearerToken,
				Client: &http.Client{
					Timeout: time.Second * 10,
				},
				Logger: log.WithFields(log.Fields{
					"sync":      "remote",
					"component": "sync",
				}),
				Cron: cron.New(),
			},
		}
		v, ok := registeredSync[name]
		if !ok {
			return nil, errors.New("no sync-provider set")
		}
		results = append(results, v)
		log.Debugf("Using %s sync-provider on %q\n", name, u)
	}

	return results, nil
}

func findEvaluator(name string) (eval.IEvaluator, error) {
	registeredEvaluators := map[string]eval.IEvaluator{
		"json": &eval.JSONEvaluator{
			Logger: log.WithFields(log.Fields{
				"evaluator": "json",
				"component": "evaluator",
			}),
		},
	}

	v, ok := registeredEvaluators[name]
	if !ok {
		return nil, errors.New("no evaluator set")
	}

	log.Debugf("Using %s evaluator\n", name)
	return v, nil
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start flagd",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Configure loggers -------------------------------------------------------
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
		// Configure service-provider impl------------------------------------------
		var serviceImpl service.IService
		serviceProvider := viper.GetString(serviceProviderFlagName)
		foundService, err := findService(serviceProvider)
		if err != nil {
			log.Errorf("Unable to find service '%s'", serviceProvider)
			return
		}
		serviceImpl = foundService

		// Configure sync-provider impl--------------------------------------------
		var syncImpl []sync.ISync
		syncProvider := viper.GetString(syncProviderFlagName)
		foundSync, err := findSync(syncProvider)
		if err != nil {
			log.Errorf("Unable to find sync '%s'", syncProvider)
			return
		}
		syncImpl = foundSync

		// Configure evaluator-provider impl------------------------------------------
		var evalImpl eval.IEvaluator
		evaluator := viper.GetString(evaluatorFlagName)
		foundEval, err := findEvaluator(evaluator)
		if err != nil {
			log.Errorf("Unable to find evaluator '%s'", evaluator)
			return
		}
		evalImpl = foundEval

		// Serve ------------------------------------------------------------------
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errc := make(chan error)
		go func() {
			errc <- func() error {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
				return fmt.Errorf("%s", <-c)
			}()
		}()

		go runtime.Start(ctx, syncImpl, serviceImpl, evalImpl, log.WithFields(log.Fields{
			"component": "runtime",
		}))
		err = <-errc
		if err != nil {
			cancel()
			log.Printf(err.Error())
		}
	},
}

func init() {
	flags := startCmd.Flags()

	// allows environment variables to use _ instead of -
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_")) // sync-provider becomes SYNC_PROVIDER
	viper.SetEnvPrefix("FLAGD")                            // port becomes FLAGD_PORT

	flags.Int32P(portFlagName, "p", 8013, "Port to listen on")
	flags.StringP(socketPathFlagName, "d", "/tmp/flagd.sock", "flagd socket path")
	flags.StringP(serviceProviderFlagName, "s", "http", "Set a service provider e.g. http or grpc")
	flags.StringP(
		syncProviderFlagName, "y", "filepath", "Set a sync provider e.g. filepath or remote",
	)
	flags.StringP(evaluatorFlagName, "e", "json", "Set an evaluator e.g. json")
	flags.StringP(serverCertPathFlagName, "c", "", "Server side tls certificate path")
	flags.StringP(serverKeyPathFlagName, "k", "", "Server side tls key path")
	flags.StringSliceP(
		uriFlagName, "f", []string{}, "Set a sync provider uri to read data from this can be a filepath or url. "+
			"Using multiple providers is supported where collisions between "+
			"flags with the same key, the later will be used.")
	flags.StringP(
		bearerTokenFlagName, "b", "", "Set a bearer token to use for remote sync")

	_ = viper.BindPFlag(portFlagName, flags.Lookup(portFlagName))
	_ = viper.BindPFlag(socketPathFlagName, flags.Lookup(socketPathFlagName))
	_ = viper.BindPFlag(serviceProviderFlagName, flags.Lookup(serviceProviderFlagName))
	_ = viper.BindPFlag(syncProviderFlagName, flags.Lookup(syncProviderFlagName))
	_ = viper.BindPFlag(evaluatorFlagName, flags.Lookup(evaluatorFlagName))
	_ = viper.BindPFlag(serverCertPathFlagName, flags.Lookup(serverCertPathFlagName))
	_ = viper.BindPFlag(serverKeyPathFlagName, flags.Lookup(serverKeyPathFlagName))
	_ = viper.BindPFlag(uriFlagName, flags.Lookup(uriFlagName))
	_ = viper.BindPFlag(bearerTokenFlagName, flags.Lookup(bearerTokenFlagName))

	_ = startCmd.MarkFlagRequired("uri")
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start flagd",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Configure loggers -------------------------------------------------------
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)

		// Build Runtime -----------------------------------------------------------
		rt, err := runtime.FromConfig(runtime.Config{
			ServiceProvider:   serviceProvider,
			ServicePort:       servicePort,
			ServiceSocketPath: socketServicePath,
			ServiceCertPath:   serverCertPath,
			ServiceKeyPath:    serverKeyPath,

			SyncProvider:    syncProvider,
			SyncURI:         uri,
			SyncBearerToken: bearerToken,

			Evaluator: evaluator,
		})
		if err != nil {
			log.Error(err)
		}

		rt.Start()
	},
}
