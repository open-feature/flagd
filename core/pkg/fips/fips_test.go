package fips

import (
	"crypto/fips140"
	"strings"
	"testing"
)

func TestStatusClassification(t *testing.T) {
	tests := []struct {
		name         string
		status       Status
		builtForFIPS bool
		certified    bool
		degraded     bool
		wantInString string
	}{
		{
			name:         "released artifact in approved mode",
			status:       Status{Enabled: true, ModuleVersion: "v1.0.0", BuildSetting: "v1.0.0-c2097c7c"},
			builtForFIPS: true,
			certified:    true,
			wantInString: "enabled (Go Cryptographic Module v1.0.0",
		},
		{
			name:         "released artifact started with GODEBUG=fips140=off",
			status:       Status{Enabled: false, ModuleVersion: "v1.0.0", BuildSetting: "v1.0.0-c2097c7c"},
			builtForFIPS: true,
			degraded:     true,
			wantInString: "GODEBUG=fips140=off at startup",
		},
		{
			name:         "local build, FIPS forced on at runtime only",
			status:       Status{Enabled: true, ModuleVersion: "latest"},
			wantInString: "uncertified",
		},
		{
			name:         "plain local build",
			status:       Status{Enabled: false, ModuleVersion: "latest"},
			wantInString: "disabled (binary not built with GOFIPS140)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.BuiltForFIPS(); got != tt.builtForFIPS {
				t.Errorf("BuiltForFIPS() = %v, want %v", got, tt.builtForFIPS)
			}
			if got := tt.status.Certified(); got != tt.certified {
				t.Errorf("Certified() = %v, want %v", got, tt.certified)
			}
			if got := tt.status.Degraded(); got != tt.degraded {
				t.Errorf("Degraded() = %v, want %v", got, tt.degraded)
			}
			if got := tt.status.String(); !strings.Contains(got, tt.wantInString) {
				t.Errorf("String() = %q, want it to contain %q", got, tt.wantInString)
			}
		})
	}
}

// TestCurrentMatchesRuntime guards the wiring between Current and the runtime,
// so the reported state can never drift from the actual FIPS mode.
func TestCurrentMatchesRuntime(t *testing.T) {
	got := Current()
	if got.Enabled != fips140.Enabled() {
		t.Errorf("Current().Enabled = %v, want %v", got.Enabled, fips140.Enabled())
	}
	if got.ModuleVersion != fips140.Version() {
		t.Errorf("Current().ModuleVersion = %q, want %q", got.ModuleVersion, fips140.Version())
	}
	// When the test binary itself is built with GOFIPS140, the two must agree.
	if got.BuiltForFIPS() && !got.Enabled {
		t.Error("built against the certified module but FIPS mode is off; flagd would refuse to start")
	}
}
