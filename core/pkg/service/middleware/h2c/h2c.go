package h2c

import (
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Middleware struct{}

func New() *Middleware {
	return &Middleware{}
}

func (m Middleware) Handler(handler http.Handler) http.Handler {
	return h2c.NewHandler(handler, &http2.Server{})
}
