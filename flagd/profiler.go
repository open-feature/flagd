//go:build profile

package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"
)

/*
Enable pprof profiler for flagd. Build controlled by the build tag "profile".
*/
func init() {
	// Go routine to server PProf
	go func() {
		server := http.Server{
			Addr:              ":6060",
			Handler:           nil,
			ReadHeaderTimeout: 3 * time.Second,
			// slowloris/slow-client DoS protection; safe for server streams
			ReadTimeout: 5 * time.Second,
		}
		server.ListenAndServe()
	}()
}
