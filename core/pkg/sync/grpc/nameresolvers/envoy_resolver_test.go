package nameresolvers

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/resolver"
)

func Test_EnvoyTargetString(t *testing.T) {
	tests := []struct {
		name        string
		mockURL     url.URL
		mockError   string
		shouldError bool
	}{
		{
			name: "Should be valid string",
			mockURL: url.URL{
				Scheme: "envoy",
				Host:   "localhost:8080",
				Path:   "/test.service",
			},
			mockError:   "",
			shouldError: false,
		},
		{
			name: "Should be valid scheme",
			mockURL: url.URL{
				Scheme: "invalid",
				Host:   "localhost:8080",
				Path:   "/test.service",
			},
			mockError:   "envoy-resolver: invalid scheme or missing host/port, target: invalid://localhost:8080/test.service",
			shouldError: true,
		},
		{
			name: "Should be valid path",
			mockURL: url.URL{
				Scheme: "envoy",
				Host:   "localhost:8080",
				Path:   "/test.service/test",
			},
			mockError:   "envoy-resolver: invalid path test.service/test",
			shouldError: true,
		},
		{
			name: "Should be valid path",
			mockURL: url.URL{
				Scheme: "envoy",
				Host:   "localhost:8080",
				Path:   "/test.service/",
			},
			mockError:   "envoy-resolver: invalid path test.service/",
			shouldError: true,
		},
		{
			name: "Hostname should not be empty",
			mockURL: url.URL{
				Scheme: "envoy",
				Host:   ":8080",
				Path:   "/test.service",
			},
			mockError:   "envoy-resolver: invalid scheme or missing host/port, target: envoy://:8080/test.service",
			shouldError: true,
		},
		{
			name: "Port should not be empty",
			mockURL: url.URL{
				Scheme: "envoy",
				Host:   "localhost",
				Path:   "/test.service",
			},
			mockError:   "envoy-resolver: invalid scheme or missing host/port, target: envoy://localhost/test.service",
			shouldError: true,
		},
		{
			name: "Hostname and Port should not be empty",
			mockURL: url.URL{
				Scheme: "envoy",
				Path:   "/test.service",
			},
			mockError:   "envoy-resolver: invalid scheme or missing host/port, target: envoy:///test.service",
			shouldError: true,
		},
	}

	for _, test := range tests {
		target := resolver.Target{URL: test.mockURL}

		isValid, err := isValidTarget(target)

		if test.shouldError {
			require.False(t, isValid, "Should not be valid")
			require.NotNilf(t, err, "Error should not be nil")
			require.Containsf(t, err.Error(), test.mockError, "Error should contains %s", test.mockError)
		} else {
			require.True(t, isValid, "Should be valid")
			require.NoErrorf(t, err, "Error should be nil")
		}
	}
}
