package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getPort(t *testing.T) {
	viper.AutomaticEnv()
	t.Cleanup(func() {
		viper.Reset()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		viper.SetEnvPrefix("FLAGD")
		bindFlags(startCmd.Flags())
	})

	tests := []struct {
		name        string
		flagName    string
		envVar      string
		envVal      string
		expected    uint16
		expectedErr string
		expectError bool
	}{
		{
			name:        "default port value when env var not set",
			flagName:    portFlagName,
			expected:    8013,
			expectError: false,
		},
		{
			name:        "default management port value when env var not set",
			flagName:    managementPortFlagName,
			expected:    8014,
			expectError: false,
		},
		{
			name:        "default sync port value when env var not set",
			flagName:    syncPortFlagName,
			expected:    8015,
			expectError: false,
		},
		{
			name:        "default ofrep port value when env var not set",
			flagName:    ofrepPortFlagName,
			expected:    8016,
			expectError: false,
		},
		{
			name:        "valid numeric port via env var",
			flagName:    portFlagName,
			envVar:      "FLAGD_PORT",
			envVal:      "9000",
			expected:    9000,
			expectError: false,
		},
		{
			name:        "leading zero decimal port resolves as base 10",
			flagName:    portFlagName,
			envVar:      "FLAGD_PORT",
			envVal:      "010000",
			expected:    10000,
			expectError: false,
		},
		{
			name:        "non-numeric port via env var returns error",
			flagName:    portFlagName,
			envVar:      "FLAGD_PORT",
			envVal:      "abc",
			expectError: true,
		},
		{
			name:        "k8s service link url via env var returns error",
			flagName:    portFlagName,
			envVar:      "FLAGD_PORT",
			envVal:      "tcp://10.0.0.1:8015",
			expectedErr: `invalid port value for 'port': "tcp://10.0.0.1:8015" (this can be caused by Kubernetes service links, set enableServiceLinks: false on the pod)`,
			expectError: true,
		},
		{
			name:        "out of range port via env var returns error",
			flagName:    portFlagName,
			envVar:      "FLAGD_PORT",
			envVal:      "70000",
			expectError: true,
		},
		{
			name:        "negative port via env var returns error",
			flagName:    portFlagName,
			envVar:      "FLAGD_PORT",
			envVal:      "-1",
			expectError: true,
		},
		{
			name:        "non-numeric management port via env var returns error",
			flagName:    managementPortFlagName,
			envVar:      "FLAGD_MANAGEMENT_PORT",
			envVal:      "invalid-port",
			expectError: true,
		},
		{
			name:        "non-numeric sync port via env var returns error",
			flagName:    syncPortFlagName,
			envVar:      "FLAGD_SYNC_PORT",
			envVal:      "invalid-port",
			expectError: true,
		},
		{
			name:        "non-numeric ofrep port via env var returns error",
			flagName:    ofrepPortFlagName,
			envVar:      "FLAGD_OFREP_PORT",
			envVal:      "invalid-port",
			expectError: true,
		},
		{
			name:        "nil flag returns 0 and nil error",
			flagName:    "non_existent_flag",
			expected:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVar != "" {
				t.Setenv(tt.envVar, tt.envVal)
			}
			got, err := getPort(tt.flagName)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid port value for '"+tt.flagName+"'")
				if tt.expectedErr != "" {
					assert.EqualError(t, err, tt.expectedErr)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
