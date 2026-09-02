package sync

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"go.uber.org/zap"
)

const (
	flagsPath = "/v1/flags"
	// only used by http.ServeContent to sniff the content type
	contentName = "flags.json"
)

// modTime records the last observed flag change. Global, since a per-selector table would be keyed
// by client input and grow unbounded.
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

// httpHandler serves the flag configuration document, not the FetchAllFlags envelope around it, so
// flagd's own HTTP sync source can consume it.
type httpHandler struct {
	store   store.IStore
	log     *logger.Logger
	modTime *modTime
}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Sampled before the read: a later stamp would strand the client on stale flags.
	lastModified := h.modTime.get()

	body, err := fetchAllFlags(r.Context(), h.store, resolveSelector(r.Header, ""))
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.Header().Set("ETag", etagOf(body))
	w.Header().Set("Vary", flagdService.FLAGD_SELECTOR_HEADER)

	// ServeContent applies RFC 9110 precedence: If-None-Match wins when both validators are sent.
	http.ServeContent(w, r, contentName, lastModified, bytes.NewReader(body))
}

// writeError maps a fetch failure onto a status. No branch echoes the expression back: it is
// unescaped client input.
func (h httpHandler) writeError(w http.ResponseWriter, err error) {
	fetchErr, ok := errors.AsType[fetchError](err)
	if !ok {
		fetchErr = fetchError{kind: fetchStoreRead, cause: err}
	}

	switch fetchErr.kind {
	case fetchSelectorMalformed:
		http.Error(w, "malformed selector", http.StatusBadRequest)
	case fetchSelectorInvalid:
		http.Error(w, "no such selector", http.StatusNotFound)
	case fetchMarshal:
		h.log.Error("error marshalling flags", zap.Error(fetchErr.cause))
		http.Error(w, "error marshalling flags", http.StatusInternalServerError)
	default:
		h.log.Error("error retrieving flags from store", zap.Error(fetchErr.cause))
		http.Error(w, "error retrieving flags from store", http.StatusInternalServerError)
	}
}

// etagOf hashes the bytes served, which is stable because encoding/json sorts map keys.
func etagOf(body []byte) string {
	sum := sha256.Sum256(body)
	return `"` + hex.EncodeToString(sum[:]) + `"`
}
