package sse

import (
	"net/http"

	"github.com/open-feature/flagd/core/pkg/store"
)

// ChannelParam carries the selector expression the client wants change notifications for, using
// the same syntax as the Flagd-Selector header; empty selects every flag. It is exported so the
// bulk handler advertises the same parameter this handler reads.
const ChannelParam = "channels"

// Handler resolves the request's selector, takes a reference on the matching subscription so the
// store watch stays alive, and streams events until the client disconnects.
func (svc *Service) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		channel := r.URL.Query().Get(ChannelParam)
		selector, err := store.NewSelector(channel)
		if err != nil {
			// not echoing the expression back: it is unescaped client input
			http.Error(w, "invalid selector in the 'channels' parameter", http.StatusBadRequest)
			return
		}

		release, err := svc.tracker.Subscribe(channel, selector)
		if err != nil {
			http.Error(w, "unable to subscribe to flag changes", http.StatusServiceUnavailable)
			return
		}
		defer release()

		svc.es.Handler(channel).ServeHTTP(w, r)
	})
}
