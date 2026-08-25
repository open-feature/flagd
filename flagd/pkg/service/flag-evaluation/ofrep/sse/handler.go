package sse

import (
	"net/http"

	"github.com/open-feature/flagd/core/pkg/store"
)

// channelPathVar is the path wildcard carrying the selector expression the client wants change
// notifications for, using the same syntax as the Flagd-Selector header. An empty channel (the
// bare stream path) selects every flag. See Service.Register for the routes it is bound to.
const channelPathVar = "channel"

// Handler resolves the request's selector, takes a reference on the matching subscription so the
// store watch stays alive, and streams events until the client disconnects.
//
// It must be mounted through Service.Register: the channel is read from the request path, so the
// route has to declare the channel wildcard.
func (svc *Service) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// already percent-decoded by the mux, so a source containing "/" arrives intact
		channel := r.PathValue(channelPathVar)
		selector, err := store.NewSelector(channel)
		if err != nil {
			// not echoing the expression back: it is unescaped client input
			http.Error(w, "invalid selector in the channel path segment", http.StatusBadRequest)
			return
		}

		release, err := svc.tracker.Subscribe(selector.ToLogString(), selector)
		if err != nil {
			http.Error(w, "unable to subscribe to flag changes", http.StatusServiceUnavailable)
			return
		}
		defer release()

		svc.es.Handler(channel).ServeHTTP(w, r)
	})
}
