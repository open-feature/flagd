package service

import (
	"net/http"
	"reflect"
	"testing"
)

const (
	headerXUserTier           = "X-User-Tier"
	headerXUserTierLowercase  = "x-user-tier"
	headerXUserEmailLowercase = "x-user-email"
)

func TestMergeContextsWithHeaders(t *testing.T) {
	type args struct {
		requestContext             map[string]any
		staticContext              map[string]any
		headers                    http.Header
		headerToContextKeyMappings map[string]string
	}

	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "empty contexts and headers",
			args: args{
				requestContext:             map[string]any{},
				staticContext:              map[string]any{},
				headers:                    http.Header{},
				headerToContextKeyMappings: map[string]string{},
			},
			want: map[string]any{},
		},
		{
			name: "request context only",
			args: args{
				requestContext:             map[string]any{"k1": "v1"},
				staticContext:              map[string]any{},
				headers:                    http.Header{},
				headerToContextKeyMappings: map[string]string{},
			},
			want: map[string]any{"k1": "v1"},
		},
		{
			name: "static context overrides request context",
			args: args{
				requestContext:             map[string]any{"k1": "v1", "k2": "v2"},
				staticContext:              map[string]any{"k2": "v22", "k3": "v3"},
				headers:                    http.Header{},
				headerToContextKeyMappings: map[string]string{},
			},
			want: map[string]any{"k1": "v1", "k2": "v22", "k3": "v3"},
		},
		{
			name: "exact case match - canonical header with canonical mapping",
			args: args{
				requestContext: map[string]any{},
				staticContext:  map[string]any{},
				headers: func() http.Header {
					h := http.Header{}
					h.Set(headerXUserTier, "premium")
					return h
				}(),
				headerToContextKeyMappings: map[string]string{headerXUserTier: "userTier"},
			},
			want: map[string]any{"userTier": "premium"},
		},
		{
			name: "case mismatch - lowercase header mapping with canonical header",
			args: args{
				requestContext: map[string]any{},
				staticContext:  map[string]any{},
				headers: func() http.Header {
					h := http.Header{}
					h.Set(headerXUserTier, "premium")
					return h
				}(),
				headerToContextKeyMappings: map[string]string{headerXUserTierLowercase: "userTier"},
			},
			want: map[string]any{"userTier": "premium"},
		},
		{
			name: "case mismatch - canonical mapping with lowercase header",
			args: args{
				requestContext: map[string]any{},
				staticContext:  map[string]any{},
				headers: func() http.Header {
					h := http.Header{}
					h.Set(headerXUserTierLowercase, "premium")
					return h
				}(),
				headerToContextKeyMappings: map[string]string{headerXUserTier: "userTier"},
			},
			want: map[string]any{"userTier": "premium"},
		},
		{
			name: "multiple headers with mixed case",
			args: args{
				requestContext: map[string]any{},
				staticContext:  map[string]any{},
				headers: func() http.Header {
					h := http.Header{}
					h.Set(headerXUserTier, "premium")
					h.Set(headerXUserEmailLowercase, "user@example.com")
					h.Set("X-Request-ID", "req-123")
					return h
				}(),
				headerToContextKeyMappings: map[string]string{
					headerXUserTierLowercase:  "userTier",
					headerXUserEmailLowercase: "userEmail",
					"x-request-id":            "requestId",
				},
			},
			want: map[string]any{
				"userTier":  "premium",
				"userEmail": "user@example.com",
				"requestId": "req-123",
			},
		},
		{
			name: "header context overrides static context",
			args: args{
				requestContext: map[string]any{"k1": "v1"},
				staticContext:  map[string]any{"k2": "v22"},
				headers: func() http.Header {
					h := http.Header{}
					h.Set("X-Override", "override-value")
					return h
				}(),
				headerToContextKeyMappings: map[string]string{"X-Override": "k2"},
			},
			want: map[string]any{"k1": "v1", "k2": "override-value"},
		},
		{
			name: "header not present - should not be in context",
			args: args{
				requestContext: map[string]any{},
				staticContext:  map[string]any{},
				headers:        http.Header{},
				headerToContextKeyMappings: map[string]string{
					"X-Missing": "missingKey",
				},
			},
			want: map[string]any{},
		},
		{
			name: "empty header value - should not be added",
			args: args{
				requestContext: map[string]any{},
				staticContext:  map[string]any{},
				headers: func() http.Header {
					h := http.Header{}
					h.Set("X-Empty", "")
					return h
				}(),
				headerToContextKeyMappings: map[string]string{
					"X-Empty": "emptyKey",
				},
			},
			want: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeContextsAndHeaders(
				tt.args.requestContext,
				tt.args.staticContext,
				tt.args.headers,
				tt.args.headerToContextKeyMappings,
			)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:  %+v\nwant: %+v", got, tt.want)
			}
		})
	}
}
