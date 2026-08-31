package sync

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	flagdService "github.com/open-feature/flagd/flagd/pkg/service"
	"go.uber.org/zap"
)

const (
	flagsPath = "/v1/flags"
	// a single path segment, percent-decoded by the mux, so a source containing "/" arrives intact
	selectorPathVar    = "selector"
	selectorQueryParam = "selector"
	// only used by http.ServeContent to sniff the content type
	contentName = "flags.json"
)

// modTime records when flagd last observed any flag change. Global rather than per-selector, since
// a per-selector table would be keyed by client-supplied strings and grow unbounded.
type modTime struct {
	unix atomic.Int64
}

func (m *modTime) set(t time.Time) {
	m.unix.Store(t.Unix())
}

// get reports the zero time until the first snapshot lands, so ServeContent omits Last-Modified.
func (m *modTime) get() time.Time {
	sec := m.unix.Load()
	if sec == 0 {
		return time.Time{}
	}
	return time.Unix(sec, 0)
}

// httpHandler is the HTTP equivalent of FetchAllFlags, sharing generateResponse with the gRPC
// handler so both wire formats stay byte-identical.
type httpHandler struct {
	store   store.IStore
	log     *logger.Logger
	modTime *modTime
}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	expression := h.selectorExpression(r)

	// Malformed is the client's mistake, an unknown filter is a miss. Neither echoes the
	// expression back: it is unescaped client input.
	if !wellFormedSelector(expression) {
		http.Error(w, "malformed selector", http.StatusBadRequest)
		return
	}

	selector, err := store.NewSelector(expression)
	if err != nil {
		http.Error(w, "no such selector", http.StatusNotFound)
		return
	}

	// Sampled before the store read: a later stamp could outrun this body, stranding the client on
	// stale flags behind a fresh Last-Modified.
	lastModified := h.modTime.get()

	flags, _, err := h.store.GetAll(r.Context(), &selector)
	if err != nil {
		h.log.Error("error retrieving flags from store", zap.Error(err))
		http.Error(w, "error retrieving flags from store", http.StatusInternalServerError)
		return
	}

	body, err := generateResponse(flags)
	if err != nil {
		h.log.Error("error marshalling flags", zap.Error(err))
		http.Error(w, "error marshalling flags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("ETag", etagOf(body))

	// ServeContent applies RFC 9110 precedence: If-None-Match wins outright when present.
	http.ServeContent(w, r, contentName, lastModified, bytes.NewReader(body))
}

// selectorExpression prefers the header, matching the precedence the gRPC and OFREP services document.
func (h httpHandler) selectorExpression(r *http.Request) string {
	if header := r.Header.Get(flagdService.FLAGD_SELECTOR_HEADER); header != "" {
		return header
	}
	if path := r.PathValue(selectorPathVar); path != "" {
		return path
	}
	return r.URL.Query().Get(selectorQueryParam)
}

// etagOf hashes the exact bytes served, which is stable because encoding/json sorts map keys.
func etagOf(body []byte) string {
	sum := sha256.Sum256(body)
	return `"` + hex.EncodeToString(sum[:]) + `"`
}

// wellFormedSelector checks shape only; unknown keys are store.NewSelector's call, and the two
// failures map to different statuses.
func wellFormedSelector(expression string) bool {
	if !utf8.ValidString(expression) {
		return false
	}
	return strings.IndexFunc(expression, unicode.IsControl) < 0
}
