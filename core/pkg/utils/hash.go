package utils

import (
	"encoding/base64"
	"encoding/json"

	"golang.org/x/crypto/sha3" //nolint:gosec
)

func GenerateSha(body []byte) string {
	hasher := sha3.New256()
	hasher.Write(canonicalize(body))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func canonicalize(body []byte) []byte {
	var parsed interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return body
	}
	canonical, err := json.Marshal(parsed)
	if err != nil {
		return body
	}
	return canonical
}
