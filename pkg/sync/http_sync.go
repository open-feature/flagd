package sync

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

type HttpSync struct {
	URI         string
	Client      *http.Client
	BearerToken string
}

func (fs *HttpSync) fetchBodyFromURL(url string) (string, error) {

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/json")

	if fs.BearerToken != "" {
		bearer := "Bearer " + fs.BearerToken
		req.Header.Set("Authorization", bearer)
	}

	resp, err := fs.Client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (fs *HttpSync) Fetch() (string, error) {
	if fs.URI == "" {
		return "", errors.New("no filepath string set")
	}

	return fs.fetchBodyFromURL(fs.URI)
}
