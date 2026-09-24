package compress

import (
	"fmt"
	"net/http"

	"github.com/klauspost/compress/gzhttp"
)

// DefaultMinSize is the smallest response worth compressing by default. See the benchmarks in this
// package: below roughly a kilobyte gzip costs more CPU than the bytes it saves, and for the very
// smallest payloads it makes the response larger.
const DefaultMinSize = 1024

type Config struct {
	// MinSize is the smallest response body, in bytes, that is compressed. Zero compresses every
	// response the content-type filter accepts, whatever its size.
	MinSize int
}

// Middleware negotiates gzip on JSON responses. Compression only happens when the client
// advertises it with Accept-Encoding, so clients that do not are unaffected.
//
// Any ETag the wrapped handler sets is deliberately left untouched: OFREP's bulk ETag is a digest
// of the uncompressed body, so a client that echoes it back in If-None-Match still gets its 304
// regardless of which encoding the original response used. gzhttp adds Vary: Accept-Encoding so
// caches keep the two encodings apart.
type Middleware struct {
	wrap func(http.Handler) http.HandlerFunc
}

func New(cfg Config) (*Middleware, error) {
	if cfg.MinSize < 0 {
		return nil, fmt.Errorf("compression min size must not be negative, got %d", cfg.MinSize)
	}

	wrap, err := gzhttp.NewWrapper(
		// gzhttp rejects a zero MinSize, so "compress everything" is expressed as a single byte.
		gzhttp.MinSize(max(cfg.MinSize, 1)),
		gzhttp.ContentTypes([]string{"application/json"}),
	)
	if err != nil {
		return nil, fmt.Errorf("building gzip middleware: %w", err)
	}

	return &Middleware{wrap: wrap}, nil
}

func (m Middleware) Handler(handler http.Handler) http.Handler {
	return m.wrap(handler)
}
