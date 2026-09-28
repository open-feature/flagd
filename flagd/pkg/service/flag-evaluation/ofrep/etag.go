package ofrep

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"iter"
	"net/http"
	"strings"
	"sync"
)

// conditionalETag tags a 200 with the written bytes and answers 304 on a match (response hashed, not config)
func conditionalETag(configEtagDiffers func(*http.Request) bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := getBodyBuffer()
		defer putBodyBuffer(body)

		rec := newResponseRecorder(w, body)
		next.ServeHTTP(rec, r)

		if rec.status != http.StatusOK {
			rec.flush()
			return
		}

		etag := rec.etag()
		w.Header().Set("ETag", etag)

		// ADR-0008 §9: a differing flagConfigEtag dictates a 200 no validator may downgrade.
		if !configEtagDiffers(r) && ifNoneMatch(r.Header.Values("If-None-Match"), etag) {
			w.Header().Del("Content-Type")
			w.Header().Del("Content-Length")
			w.WriteHeader(http.StatusNotModified)
			return
		}

		rec.flush()
	})
}

// bodyBuffers recycles response buffers; consecutive bulk evaluations are close in size
var bodyBuffers = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

// getBodyBuffer lends out an empty buffer.
func getBodyBuffer() *bytes.Buffer {
	buf := bodyBuffers.Get().(*bytes.Buffer)
	buf.Reset()

	return buf
}

// putBodyBuffer returns a buffer whatever it grew to.
func putBodyBuffer(buf *bytes.Buffer) {
	bodyBuffers.Put(buf)
}

type responseRecorder struct {
	http.ResponseWriter

	status  int
	body    *bytes.Buffer
	digest  hash.Hash
	written bool
}

// newResponseRecorder records into body, hashing as it goes.
func newResponseRecorder(w http.ResponseWriter, body *bytes.Buffer) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		status:         http.StatusOK,
		body:           body,
		digest:         sha256.New(),
	}
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
	_, _ = rec.digest.Write(b) // never errors

	return rec.body.Write(b)
}

// etag returns a weak validator; gzip and identity share it since they cannot share a strong tag (RFC 9110 8.8.1)
func (rec *responseRecorder) etag() string {
	return `W/"` + hex.EncodeToString(rec.digest.Sum(nil)) + `"`
}

func (rec *responseRecorder) flush() {
	rec.ResponseWriter.WriteHeader(rec.status)
	if rec.body.Len() > 0 {
		_, _ = rec.ResponseWriter.Write(rec.body.Bytes())
	}
}

// ifNoneMatch reports whether any If-None-Match value selects the representation tagged with etag (RFC 9110 13.1.2)
func ifNoneMatch(fields []string, etag string) bool {
	opaque := strings.TrimPrefix(etag, "W/")

	for _, field := range fields {
		for candidate := range splitETagList(field) {
			if candidate == "*" {
				return true
			}

			// a lenient client may drop the quotes the generated tag carries
			candidate = strings.TrimPrefix(candidate, "W/")
			if candidate == opaque || `"`+candidate+`"` == opaque {
				return true
			}
		}
	}

	return false
}

// splitETagList splits a list of entity tags on commas, tracking quoting since a quoted tag may contain one
func splitETagList(list string) iter.Seq[string] {
	return func(yield func(string) bool) {
		var (
			start  int
			quoted bool
		)
		flush := func(end int) bool {
			if tag := strings.TrimSpace(list[start:end]); tag != "" {
				if !yield(tag) {
					return false
				}
			}
			return true
		}
		for i := 0; i < len(list); i++ {
			switch c := list[i]; {
			case c == '"':
				quoted = !quoted
			case c == ',' && !quoted:
				if !flush(i) {
					return
				}
				start = i + 1
			}
		}
		flush(len(list))
	}
}
