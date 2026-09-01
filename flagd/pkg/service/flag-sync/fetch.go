package sync

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/open-feature/flagd/core/pkg/store"
	flagdService "github.com/open-feature/flagd/flagd/pkg/service"
)

// resolveSelector prefers the Flagd-Selector header, as the other flagd services do.
func resolveSelector(header http.Header, fallback string) string {
	if value := header.Get(flagdService.FLAGD_SELECTOR_HEADER); value != "" {
		return value
	}
	return fallback
}

type fetchErrorKind int

const (
	fetchSelectorMalformed fetchErrorKind = iota
	fetchSelectorInvalid
	fetchStoreRead
	fetchMarshal
)

// fetchError classifies a failure so each transport maps it onto its own status vocabulary.
type fetchError struct {
	kind  fetchErrorKind
	cause error
}

func (e fetchError) Error() string { return e.cause.Error() }

func (e fetchError) Unwrap() error { return e.cause }

// newSelector validates shape before keys, since the two failures map to different statuses.
func newSelector(expression string) (store.Selector, error) {
	if !wellFormedSelector(expression) {
		return store.Selector{}, fetchError{kind: fetchSelectorMalformed, cause: errors.New("malformed selector")}
	}

	selector, err := store.NewSelector(expression)
	if err != nil {
		return store.Selector{}, fetchError{kind: fetchSelectorInvalid, cause: err}
	}

	return selector, nil
}

// fetchAllFlags is the transport-agnostic body of FetchAllFlags.
func fetchAllFlags(ctx context.Context, s store.IStore, expression string) ([]byte, error) {
	selector, err := newSelector(expression)
	if err != nil {
		return nil, err
	}

	flags, _, err := s.GetAll(ctx, &selector)
	if err != nil {
		return nil, fetchError{kind: fetchStoreRead, cause: err}
	}

	body, err := generateResponse(flags)
	if err != nil {
		return nil, fetchError{kind: fetchMarshal, cause: err}
	}

	return body, nil
}

// wellFormedSelector checks shape only; unknown keys are store.NewSelector's call.
func wellFormedSelector(expression string) bool {
	if !utf8.ValidString(expression) {
		return false
	}
	return strings.IndexFunc(expression, unicode.IsControl) < 0
}
