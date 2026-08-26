package ofrep

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"net/http"
	"strings"
	"sync"
)

// conditionalETag tags a 200 with a strong ETag over the bytes the handler wrote, and answers 304
// when the client's If-None-Match matches. Hashing the response rather than the flag configuration
// keeps the validator correct for context-dependent values.
func conditionalETag(next http.Handler) http.Handler {
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
		w.Header().Set("ETag", quoteETag(etag))

		if ifNoneMatch(r.Header.Values("If-None-Match"), etag) {
			w.Header().Del("Content-Type")
			w.Header().Del("Content-Length")
			w.WriteHeader(http.StatusNotModified)
			return
		}

		rec.flush()
	})
}

// bodyBuffers recycles response buffers: consecutive bulk evaluations are close in size, so a
// recycled buffer already holds the capacity the next one needs.
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

func (rec *responseRecorder) etag() string {
	return hex.EncodeToString(rec.digest.Sum(nil))
}

func (rec *responseRecorder) flush() {
	rec.ResponseWriter.WriteHeader(rec.status)
	if rec.body.Len() > 0 {
		_, _ = rec.ResponseWriter.Write(rec.body.Bytes())
	}
}

// ifNoneMatch reports whether any of the client's If-None-Match field values selects the
// representation tagged with etag. Per RFC 9110 §13.1.2 the field is either "*" or a list of
// entity tags, repeated fields extend that list, and the comparison is weak.
//
// net/http parses this too, but scanETag and checkIfNoneMatch are unexported, and the one
// exported way in, http.ServeContent, answers a failed precondition on a non-GET method with 412
// rather than the 304 OFREP asks for on this POST route.
//
// The scan is more forgiving than net/http's, which abandons the whole field on the first entry
// that is not a syntactically valid entity tag.
func ifNoneMatch(fields []string, etag string) bool {
	for _, field := range fields {
		for _, candidate := range splitETagList(field) {
			if candidate == "*" || normalizeETag(candidate) == etag {
				return true
			}
		}
	}

	return false
}

// splitETagList splits a list of entity tags on its commas. A comma inside a quoted tag belongs
// to the tag, so the split tracks quoting instead of reaching for strings.Split.
func splitETagList(list string) []string {
	var (
		tags    []string
		current strings.Builder
		quoted  bool
	)

	flush := func() {
		if tag := strings.TrimSpace(current.String()); tag != "" {
			tags = append(tags, tag)
		}
		current.Reset()
	}

	for i := 0; i < len(list); i++ {
		switch c := list[i]; {
		case c == '"':
			quoted = !quoted
			current.WriteByte(c)
		case c == ',' && !quoted:
			flush()
		default:
			current.WriteByte(c)
		}
	}
	flush()

	return tags
}

// normalizeETag reduces an entity tag to the opaque value the weak comparison function looks at:
// neither the quotes nor a weakness prefix carries meaning there. Dropping the quotes also lets
// the unquoted form a lenient client may send compare equal.
func normalizeETag(etag string) string {
	return strings.Trim(strings.TrimPrefix(etag, "W/"), `"`)
}

func quoteETag(etag string) string {
	return `"` + etag + `"`
}
