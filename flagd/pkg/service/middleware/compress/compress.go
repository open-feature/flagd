package compress

import (
	"fmt"
	"net/http"

	"github.com/klauspost/compress/gzhttp"
)

// Middleware gzips JSON responses. Compression only happens when the client advertises it with
// Accept-Encoding, so clients that do not are unaffected, and gzhttp adds Vary: Accept-Encoding so
// caches keep the two encodings apart.
//
// Every response is compressed regardless of size. Gzip costs a near-fixed few microseconds, which
// a small body does not repay, but a size threshold would only trade that for a configuration knob
// and an encoding the handler below cannot predict.
type Middleware struct {
	wrap func(http.Handler) http.HandlerFunc
}

// New builds the gzip middleware.
func New() (*Middleware, error) {
	wrap, err := gzhttp.NewWrapper(
		// 1, not 0: at 0 gzhttp starts the encoder before it has seen any body, so a 304 would go
		// out carrying Content-Encoding and gzip framing. At 1 an empty body stays plain and
		// everything else is compressed, whatever its size.
		gzhttp.MinSize(1),
		gzhttp.ContentTypes([]string{"application/json"}),
	)
	if err != nil {
		return nil, fmt.Errorf("building gzip middleware: %w", err)
	}

	return &Middleware{wrap: wrap}, nil
}

// Handler wraps the next handler with gzip negotiation.
func (m Middleware) Handler(handler http.Handler) http.Handler {
	return m.wrap(handler)
}
