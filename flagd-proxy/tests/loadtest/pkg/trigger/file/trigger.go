package trigger

import (
	"fmt"
	"os"
)

type FilePathTrigger struct {
	config FilePathTriggerConfig
}

type FilePathTriggerConfig struct {
	StartFile  string `json:"startFile"`
	EndFile    string `json:"endFile"`
	TargetFile string `json:"targetFile"`
}

func NewFilePathTrigger(config FilePathTriggerConfig) *FilePathTrigger {
	return &FilePathTrigger{
		config: config,
	}
}

func (f *FilePathTrigger) Setup() error {
	dat, err := os.ReadFile(f.config.StartFile)
	if err != nil {
		return fmt.Errorf("unable to read start file at %s: %w", f.config.StartFile, err)
	}
	if err = os.WriteFile(f.config.TargetFile, dat, 0o600); err != nil {
		return fmt.Errorf("unable to write start file: %w", err)
	}
	return nil
}

func (f *FilePathTrigger) Update() error {
	dat, err := os.ReadFile(f.config.EndFile)
	if err != nil {
		return fmt.Errorf("unable to read end file at %s: %w", f.config.EndFile, err)
	}
	if err = os.WriteFile(f.config.TargetFile, dat, 0o600); err != nil {
		return fmt.Errorf("unable to write end file: %w", err)
	}
	return nil
}
