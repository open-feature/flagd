package middleware

import (
	"fmt"
	"net/http"
)

type Middleware interface {
	Handle(handler http.Handler) http.Handler
}

type Logger struct{}

func (Logger) Handle(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("log")
		handler.ServeHTTP(writer, request)
	})
}
