package compress

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/klauspost/compress/gzhttp"
)

// DefaultMinSize is the smallest response worth compressing by default. See the benchmarks in this
// package: gzip costs a near-fixed few microseconds per response whatever its size, and below
// roughly a kilobyte that buys almost nothing. A single-flag evaluation is around 150 bytes and
// barely shrinks at all, while a bulk response of ten flags already compresses several times over.
const DefaultMinSize = 1024

// ETagSuffix marks the gzip-encoded representation of a response. Content coding is part of the
// representation, so per RFC 9110 8.8.1 the compressed and uncompressed forms of a body must not
// share one strong validator. Appending this to the ETag on compressed responses keeps the two
// apart; a handler that answers conditional requests has to strip it back off with TrimETagSuffix.
const ETagSuffix = "-gzip"

type Config struct {
	// MinSize is the smallest response body, in bytes, that is compressed. Zero compresses every
	// response the content-type filter accepts, whatever its size.
	MinSize int
}

// Middleware negotiates gzip on JSON responses. Compression only happens when the client
// advertises it with Accept-Encoding, so clients that do not are unaffected. gzhttp adds
// Vary: Accept-Encoding so caches keep the two encodings apart.
type Middleware struct {
	wrap func(http.Handler) http.HandlerFunc
}

func New(cfg Config) (*Middleware, error) {
	if cfg.MinSize < 0 {
		return nil, fmt.Errorf("compression min size must not be negative, got %d", cfg.MinSize)
	}

	wrap, err := gzhttp.NewWrapper(
		gzhttp.MinSize(cfg.MinSize),
		gzhttp.ContentTypes([]string{"application/json"}),
		gzhttp.SuffixETag(ETagSuffix),
	)
	if err != nil {
		return nil, fmt.Errorf("building gzip middleware: %w", err)
	}

	return &Middleware{wrap: wrap}, nil
}

func (m Middleware) Handler(handler http.Handler) http.Handler {
	return m.wrap(handler)
}

// TrimETagSuffix removes the gzip marker from an entity tag, yielding the tag of the uncompressed
// representation. Tags that carry no marker are returned unchanged.
func TrimETagSuffix(etag string) string {
	if trimmed, ok := strings.CutSuffix(etag, ETagSuffix+`"`); ok {
		return trimmed + `"`
	}

	return strings.TrimSuffix(etag, ETagSuffix)
}

// AcceptsGzip reports whether the client is willing to receive a gzip-encoded response, and so
// whether a gzip-tagged validator it sends back can describe the representation it would get now.
// Deliberately lenient: it only gates a validator this server issued in the first place.
func AcceptsGzip(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("Accept-Encoding")), "gzip")
}
