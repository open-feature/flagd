package compress

import (
	"fmt"
	"net/http"

	"github.com/klauspost/compress/gzhttp"
)

// minSize is the smallest response worth compressing. A single-flag OFREP evaluation is a couple
// of hundred bytes, where the gzip framing costs more than it saves; bulk evaluations are far
// larger and compress well.
const minSize = 1024

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

func New() (*Middleware, error) {
	wrap, err := gzhttp.NewWrapper(
		gzhttp.MinSize(minSize),
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
