package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"

	"golang.org/x/crypto/sha3" //nolint:gosec
)

func GenerateSha(body []byte) string {
	hasher := sha3.New256()
	hasher.Write(canonicalize(body))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func canonicalize(body []byte) []byte {
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()

	var parsed interface{}
	if err := dec.Decode(&parsed); err != nil {
		return body
	}
	// check for leftover garbage after valid json
	var extra json.RawMessage
	if dec.Decode(&extra) != io.EOF {
		return body
	}

	canonical, err := json.Marshal(parsed)
	if err != nil {
		return body
	}
	return canonical
}
