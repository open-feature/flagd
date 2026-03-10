package service

import (
	"net/http"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestSelectorExpressionFromHTTPHeaders(t *testing.T) {
	headers := http.Header{}
	headers.Add(FLAG_SELECTOR_HEADER, "A")
	headers.Add(FLAG_SELECTOR_HEADER, "B")

	got := SelectorExpressionFromHTTPHeaders(headers)
	if got != "A,B" {
		t.Fatalf("expected A,B, got %s", got)
	}
}

func TestSelectorExpressionFromGRPCMetadata(t *testing.T) {
	md := metadata.New(map[string]string{
		FLAG_SELECTOR_HEADER: "A,C,B",
	})

	got := SelectorExpressionFromGRPCMetadata(md)
	if got != "A,C,B" {
		t.Fatalf("expected A,C,B, got %s", got)
	}
}
