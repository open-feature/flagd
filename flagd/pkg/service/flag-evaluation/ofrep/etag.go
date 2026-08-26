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

		if ifNoneMatch(r.Header.Values("If-None-Match"), etag) {
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

// ifNoneMatch reports whether any of the client's If-None-Match field values selects the
// representation tagged with etag. Per RFC 9110 §13.1.2 the field is either "*" or a list of
// entity tags, repeated fields extend that list, and the comparison is weak.
//
// net/http parses this too, but scanETag and checkIfNoneMatch are unexported, and the one
// exported way in, http.ServeContent, answers a failed precondition on a non-GET method with 412
// rather than the 304 OFREP asks for on this POST route.
//
// The scan is more forgiving than net/http's, which abandons the whole field on the first entry
// that is not a syntactically valid entity tag. A malformed entry is compared as an opaque token
// instead, so neither it nor an unquoted tag hides a valid tag listed beside it. Erring this way
// only ever costs a client the bandwidth of a body it already holds.
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
