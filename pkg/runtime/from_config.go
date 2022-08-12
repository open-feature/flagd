package runtime

import (
	"github.com/open-feature/flagd/pkg/sync"
	log "github.com/sirupsen/logrus"
)

func RuntimeFromConfig(config RuntimeConfig) (*Runtime, error) {
	rt := Runtime{
		config: config,
		Logger: log.WithFields(log.Fields{
			"component": "runtime",
		}),
		syncNotifier: make(chan sync.INotify),
	}
	if err := rt.SetEvaluator(); err != nil {
		return nil, err
	}
	if err := rt.SetService(); err != nil {
		return nil, err
	}
	if err := rt.SetSyncImpl(); err != nil {
		return nil, err
	}
	return &rt, nil
}
