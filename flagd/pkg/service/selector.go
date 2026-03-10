package service
 
import (
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

func SelectorExpressionFromHTTPHeaders(headers http.Header) string {
	if selectors := strings.TrimSpace(strings.Join(headers.Values(FLAG_SELECTOR_HEADER), ",")); selectors != "" {
		return selectors
	}
	return strings.TrimSpace(strings.Join(headers.Values(FLAGD_SELECTOR_HEADER), ","))
}

func SelectorExpressionFromGRPCMetadata(md metadata.MD) string {
	if selectors := strings.TrimSpace(strings.Join(md.Get(strings.ToLower(FLAG_SELECTOR_HEADER)), ",")); selectors != "" {
		return selectors
	}
	if selectors := strings.TrimSpace(strings.Join(md.Get(strings.ToLower(FLAGD_SELECTOR_HEADER)), ",")); selectors != "" {
		return selectors
	}
	if selectors := strings.TrimSpace(strings.Join(md.Get(FLAG_SELECTOR_HEADER), ",")); selectors != "" {
		return selectors
	}
	return strings.TrimSpace(strings.Join(md.Get(FLAGD_SELECTOR_HEADER), ","))
}
