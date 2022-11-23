//go:build profile

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

// Init PProf server
func init() {
	// Go routine to server PProf
	go func() {
		server := http.Server{Addr: "localhost:6060", Handler: nil}
		err := server.ListenAndServe()

		if err != nil {
			fmt.Printf("Server start : %s", err)
		}
	}()
}
