package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/open-feature/flagd/flagd-proxy/tests/loadtest/pkg/handler"
	trigger "github.com/open-feature/flagd/flagd-proxy/tests/loadtest/pkg/trigger/file"
)

type TriggerType string

const (
	FilepathTrigger TriggerType = "filepath"

	defaultHost                    = "localhost"
	defaultPort             uint16 = 8080
	defaultStartFile               = "./config/start-spec.json"
	defaultEndFile                 = "./config/end-spec.json"
	defaultTargetFile              = "./target.json"
	defaultTargetFileSource        = "./target.json"
	defaultOutTarget               = "./profiling-results.json"
)

var defaultTests = []handler.TestConfig{
	{
		Watchers: 1,
		Repeats:  5,
		Delay:    time.Second * 1,
	},
	{
		Watchers: 10,
		Repeats:  5,
		Delay:    time.Second * 1,
	},
	{
		Watchers: 100,
		Repeats:  5,
		Delay:    time.Second * 1,
	},
	{
		Watchers: 1000,
		Repeats:  5,
		Delay:    time.Second * 1,
	},
	{
		Watchers: 10000,
		Repeats:  5,
		Delay:    time.Second * 1,
	},
}

type Config struct {
	TriggerType       TriggerType                   `json:"triggerType"`
	FileTriggerConfig trigger.FilePathTriggerConfig `json:"fileTriggerConfig"`
	HandlerConfig     handler.Config                `json:"handlerConfig"`
	Tests             []handler.TestConfig
}

func NewConfig(filepath string) (*Config, error) {
	config := &Config{
		TriggerType: FilepathTrigger,
		FileTriggerConfig: trigger.FilePathTriggerConfig{
			StartFile:  defaultStartFile,
			EndFile:    defaultEndFile,
			TargetFile: defaultTargetFile,
		},
		HandlerConfig: handler.Config{
			FilePath: defaultTargetFileSource,
			Host:     defaultHost,
			Port:     defaultPort,
			OutFile:  defaultOutTarget,
		},
		Tests: defaultTests,
	}
	if filepath != "" {
		b, err := os.ReadFile(filepath)
		if err != nil {
			return nil, fmt.Errorf("unable to read config file %s: %w", filepath, err)
		}
		if err := json.Unmarshal(b, config); err != nil {
			return nil, fmt.Errorf("unable to unmarshal config: %w", err)
		}
	}
	return config, nil
}
