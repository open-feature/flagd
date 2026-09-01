package runtime

import (
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// An operator who set one has to be told rather than left assuming it took effect.
//
//nolint:staticcheck // SA1019 the no-op fields are what is under test
func TestWarnOnNoOpConfig(t *testing.T) {
	const snippet = "are no-ops"

	tests := []struct {
		name     string
		config   Config
		wantWarn bool
	}{
		{
			name: "defaults stay silent",
			config: Config{
				KeepAliveMinTime:             KeepAliveMinTimeDefault,
				KeepAlivePermitWithoutStream: KeepAlivePermitWithoutStreamDefault,
			},
		},
		{
			name: "changed min time warns",
			config: Config{
				KeepAliveMinTime:             10 * time.Second,
				KeepAlivePermitWithoutStream: KeepAlivePermitWithoutStreamDefault,
			},
			wantWarn: true,
		},
		{
			name: "changed permit without stream warns",
			config: Config{
				KeepAliveMinTime:             KeepAliveMinTimeDefault,
				KeepAlivePermitWithoutStream: false,
			},
			wantWarn: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			core, logs := observer.New(zap.WarnLevel)

			warnOnNoOpConfig(logger.NewLogger(zap.New(core), false), test.config)

			warned := logs.FilterMessageSnippet(snippet).Len() > 0
			assert.Equal(t, test.wantWarn, warned)
		})
	}
}
