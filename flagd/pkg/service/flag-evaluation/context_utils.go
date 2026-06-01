package service

import (
	"net/http"

	"connectrpc.com/connect"
	flagdService "github.com/open-feature/flagd/flagd/pkg/service"
	"github.com/open-feature/flagd/core/pkg/store"
)

func selectorFromHeader(header http.Header) (store.Selector, error) {
	expr := header.Get(flagdService.FLAGD_SELECTOR_HEADER)
	s, err := store.NewSelector(expr)
	if err != nil {
		return store.Selector{}, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return s, nil
}

// MergeContextsAndHeaders merges evaluation contexts with static context values and header-based context.
// highest priority > header-context-from-cli > static-context-from-cli > request-context > lowest priority
// Header names are matched case-insensitively according to HTTP specification.
func MergeContextsAndHeaders(
	requestContext map[string]any,
	staticContext map[string]any,
	headers http.Header,
	headerToContextKeyMappings map[string]string,
) map[string]any {
	merged := make(map[string]any)

	// request-body/client context first (lowest priority)
	for k, v := range requestContext {
		merged[k] = v
	}

	// static/config context (overrides request context)
	for k, v := range staticContext {
		merged[k] = v
	}

	// header-derived context (highest priority) we use .Get which is case-insensitive
	for headerName, contextKey := range headerToContextKeyMappings {
		if value := headers.Get(headerName); value != "" {
			merged[contextKey] = value
		}
	}

	return merged
}
