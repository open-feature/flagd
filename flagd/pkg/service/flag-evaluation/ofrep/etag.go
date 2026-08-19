package ofrep

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

// conditionalETag buffers the bytes the wrapped handler writes, tags them with a strong ETag and
// answers 304 when the client's If-None-Match names that representation. Hashing the octets means
// the validator covers context-dependent values too, where a config-derived one would let a
// client that changed a context attribute keep evaluations resolved against the old one.
//
// Only 200s are tagged, and the tag needs the response, so a 304 saves bandwidth, not work.
func conditionalETag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		if rec.status != http.StatusOK {
			rec.flush()
			return
		}

		etag := bodyETag(rec.body.Bytes())
		w.Header().Set("ETag", quoteETag(etag))

		if clientETag := r.Header.Get("If-None-Match"); clientETag != "" && normalizeETag(clientETag) == etag {
			// A 304 carries no body, so drop the headers that describe one.
			w.Header().Del("Content-Type")
			w.Header().Del("Content-Length")
			w.WriteHeader(http.StatusNotModified)
			return
		}

		rec.flush()
	})
}

// responseRecorder holds back a handler's status and body so they can be hashed before anything
// reaches the client. Headers pass straight through to the real writer.
type responseRecorder struct {
	http.ResponseWriter

	status  int
	body    bytes.Buffer
	written bool
}

func (rec *responseRecorder) WriteHeader(status int) {
	if rec.written {
		return
	}
	rec.status = status
	rec.written = true
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	if !rec.written {
		rec.WriteHeader(http.StatusOK)
	}
	return rec.body.Write(b)
}

func (rec *responseRecorder) flush() {
	rec.ResponseWriter.WriteHeader(rec.status)
	if rec.body.Len() > 0 {
		_, _ = rec.ResponseWriter.Write(rec.body.Bytes())
	}
}

// bodyETag hashes the response. Encoding is deterministic (struct field order, sorted map
// keys), so an unchanged representation keeps the same tag across requests.
func bodyETag(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

// normalizeETag strips optional quotes so quoted and unquoted forms compare equal.
func normalizeETag(etag string) string {
	return strings.Trim(etag, `"`)
}

func quoteETag(etag string) string {
	return `"` + etag + `"`
}
