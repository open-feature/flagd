package utils

import (
	"bytes"
	"crypto/sha3"
	"encoding/base64"
	"encoding/json"
	"hash"
	"io"
)

// GenerateSha returns a SHA3-256 digest of the canonicalized body.
func GenerateSha(body []byte) string {
	// hash.Hash rather than *sha3.SHA3: its Write never returns an error.
	var hasher hash.Hash = sha3.New256()
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
