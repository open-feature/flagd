package sync

import (
	"errors"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

type FilePathSync struct {
	URI string
}

func (fs *FilePathSync) Fetch() (string, error) {
	if fs.URI == "" {
		return "", errors.New("no filepath string set")
	}

	rawFile, err := ioutil.ReadFile(fs.URI)
	if err != nil {
		return "", err
	}
	log.Debugf("Fetched file: ", fs.URI)
	return string(rawFile), nil
}
