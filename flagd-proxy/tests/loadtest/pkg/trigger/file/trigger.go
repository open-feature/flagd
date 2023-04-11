package trigger

import "os"

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
		return err
	}
	return os.WriteFile(f.config.TargetFile, dat, 0644)
}

func (f *FilePathTrigger) Update() error {
	dat, err := os.ReadFile(f.config.EndFile)
	if err != nil {
		return err
	}
	return os.WriteFile(f.config.TargetFile, dat, 0644)
}

func Cleanup(filename string) error {
	return os.Remove(filename)
}
