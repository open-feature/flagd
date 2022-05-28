package sync

import (
	"errors"
	"io/ioutil"
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
	return string(rawFile), nil
}
