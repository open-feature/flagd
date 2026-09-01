package sync

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/logger"
	flagdService "github.com/open-feature/flagd/flagd/pkg/service"
	"github.com/stretchr/testify/require"
)

const selectorHeaderKey = flagdService.FLAGD_SELECTOR_HEADER

func TestResolveSelector(t *testing.T) {
	tests := []struct {
		name     string
		header   http.Header
		fallback string
		want     string
	}{
		{name: "header wins", header: selectorHeader("source=a"), fallback: "source=b", want: "source=a"},
		{name: "fallback when header absent", header: http.Header{}, fallback: "source=b", want: "source=b"},
		{name: "fallback when header empty", header: selectorHeader(""), fallback: "source=b", want: "source=b"},
		{name: "nil header", header: nil, fallback: "source=b", want: "source=b"},
		{name: "neither", header: nil, fallback: "", want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, resolveSelector(test.header, test.fallback))
		})
	}
}

func TestNewSelector_ErrorKinds(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantKind   fetchErrorKind
	}{
		{name: "control character", expression: "source=a\x00b", wantKind: fetchSelectorMalformed},
		{name: "invalid utf8", expression: "source=\xc3\x28", wantKind: fetchSelectorMalformed},
		{name: "unknown key", expression: "invalidKey=val", wantKind: fetchSelectorInvalid},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := newSelector(test.expression)

			var fetchErr fetchError
			require.ErrorAs(t, err, &fetchErr)
			require.Equal(t, test.wantKind, fetchErr.kind)
		})
	}

	t.Run("empty selector is valid", func(t *testing.T) {
		_, err := newSelector("")
		require.NoError(t, err)
	})
}

func TestFetchAllFlags(t *testing.T) {
	flagStore, _ := getSimpleFlagStore(t)

	t.Run("returns the whole configuration", func(t *testing.T) {
		body, err := fetchAllFlags(context.Background(), flagStore, "")
		require.NoError(t, err)
		require.Contains(t, string(body), "flagA")
		require.Contains(t, string(body), "flagB")
	})

	t.Run("classifies a bad selector", func(t *testing.T) {
		_, err := fetchAllFlags(context.Background(), flagStore, "invalidKey=val")

		var fetchErr fetchError
		require.ErrorAs(t, err, &fetchErr)
		require.Equal(t, fetchSelectorInvalid, fetchErr.kind)
	})
}

// Both selector failures share a code, and only they surface their cause to the client.
func TestConnectError(t *testing.T) {
	h := syncHandler{log: logger.NewLogger(nil, false)}

	tests := []struct {
		name        string
		err         error
		wantCode    connect.Code
		wantMessage string
	}{
		{
			name:        "malformed selector",
			err:         fetchError{kind: fetchSelectorMalformed, cause: errors.New("malformed selector")},
			wantCode:    connect.CodeInvalidArgument,
			wantMessage: "malformed selector",
		},
		{
			name:        "invalid selector",
			err:         fetchError{kind: fetchSelectorInvalid, cause: errors.New("invalid selector key")},
			wantCode:    connect.CodeInvalidArgument,
			wantMessage: "invalid selector key",
		},
		{
			name:        "store read hides the cause",
			err:         fetchError{kind: fetchStoreRead, cause: errors.New("disk on fire")},
			wantCode:    connect.CodeInternal,
			wantMessage: "error retrieving flags from store",
		},
		{
			name:        "marshal hides the cause",
			err:         fetchError{kind: fetchMarshal, cause: errors.New("cycle in value")},
			wantCode:    connect.CodeDataLoss,
			wantMessage: "error marshalling flags",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requireConnectError(t, h.connectError(test.err), test.wantCode, test.wantMessage)
		})
	}

	t.Run("unclassified errors pass through", func(t *testing.T) {
		sentinel := errors.New("not from a fetch")
		require.Equal(t, sentinel, h.connectError(sentinel))
	})
}

func BenchmarkFetchAllFlags(b *testing.B) {
	flagStore, _ := getSimpleFlagStore(b)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := fetchAllFlags(ctx, flagStore, ""); err != nil {
			b.Fatal(err)
		}
	}
}

func selectorHeader(value string) http.Header {
	header := http.Header{}
	header.Set(selectorHeaderKey, value)
	return header
}
