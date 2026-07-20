package h2c

import (
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c" //nolint:staticcheck // deprecated package, see Handler below
)

type Middleware struct{}

func New() *Middleware {
	return &Middleware{}
}

func (m Middleware) Handler(handler http.Handler) http.Handler {
	// h2c.NewHandler is deprecated in favor of setting http.Server.Protocols,
	// but that requires refactoring server construction; the handler still
	// works correctly, so we keep using it for now.
	//nolint:staticcheck
	return h2c.NewHandler(handler, &http2.Server{})
}
