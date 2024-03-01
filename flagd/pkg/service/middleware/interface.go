package middleware

import (
	"net/http"
)

type IMiddleware interface {
	Handler(handler http.Handler) http.Handler
}
