package ofrep

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

// conditionalETag tags a 200 with a strong ETag over the bytes the handler wrote, and answers 304
// when the client's If-None-Match matches. Hashing the response rather than the flag configuration
// keeps the validator correct for context-dependent values; the trade is that producing the tag
// needs the response, so a 304 saves bandwidth, not work.
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
			w.Header().Del("Content-Type")
			w.Header().Del("Content-Length")
			w.WriteHeader(http.StatusNotModified)
			return
		}

		rec.flush()
	})
}

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

func bodyETag(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

// normalizeETag lets quoted and unquoted forms compare equal.
func normalizeETag(etag string) string {
	return strings.Trim(etag, `"`)
}

func quoteETag(etag string) string {
	return `"` + etag + `"`
}
