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
		wantInString string
	}{
		{
			name:         "released FIPS artifact in approved mode",
			status:       Status{Enabled: true, ModuleVersion: "v1.0.0", BuildSetting: CertifiedModuleSnapshot},
			builtForFIPS: true,
			certified:    true,
			wantInString: "enabled (Go Cryptographic Module v1.0.0",
		},
		{
			name:         "FIPS artifact started with GODEBUG=fips140=off",
			status:       Status{Enabled: false, ModuleVersion: "v1.0.0", BuildSetting: CertifiedModuleSnapshot},
			builtForFIPS: true,
			wantInString: "GODEBUG=fips140=off at startup",
		},
		{
			name:         "a different frozen module snapshot is not the certified one",
			status:       Status{Enabled: true, ModuleVersion: "v1.0.0", BuildSetting: "v1.0.0-deadbeef"},
			wantInString: "uncertified",
		},
		{
			name:         "standard build with FIPS forced on at runtime",
			status:       Status{Enabled: true, ModuleVersion: "latest"},
			wantInString: "uncertified",
		},
		{
			name:         "standard build",
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
}

// TestCheckMatchesVariant runs under both build variants: the standard build
// must never refuse, and the FIPS build must refuse exactly when the certified
// module is not active.
func TestCheckMatchesVariant(t *testing.T) {
	err := Check()
	if !RequireCertified {
		if Variant != "standard" {
			t.Errorf("Variant = %q, want %q", Variant, "standard")
		}
		if err != nil {
			t.Errorf("standard build: Check() = %v, want nil", err)
		}
		return
	}
	if Variant != "fips" {
		t.Errorf("Variant = %q, want %q", Variant, "fips")
	}
	if certified := Current().Certified(); certified != (err == nil) {
		t.Errorf("Check() = %v but Certified() = %v; they must agree", err, certified)
	}
}
