package cors

import (
	"net/http"

	"github.com/rs/cors"
)

type Middleware struct {
	cors *cors.Cors
}

func New(allowedOrigins []string) *Middleware {
	return &Middleware{
		cors: cors.New(cors.Options{
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedOrigins: allowedOrigins,
			AllowedHeaders: []string{"*"},
			ExposedHeaders: []string{
				// Content-Type is in the default safelist.
				"Accept",
				"Accept-Encoding",
				"Accept-Post",
				"Connect-Accept-Encoding",
				"Connect-Content-Encoding",
				"Content-Encoding",
				"Grpc-Accept-Encoding",
				"Grpc-Encoding",
				"Grpc-Message",
				"Grpc-Status",
				"Grpc-Status-Details-Bin",
			},
		}),
	}
}

func (c Middleware) Handler(handler http.Handler) http.Handler {
	return c.cors.Handler(handler)
}
