package middleware

import (
	"net/http"
)

type Middleware interface {
	Handle(handler http.Handler) http.Handler
}
