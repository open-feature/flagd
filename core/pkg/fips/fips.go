// Package fips reports whether this binary is running the CMVP-certified Go
// Cryptographic Module in FIPS 140-3 approved mode.
//
// flagd is released in two variants. The standard build uses the ordinary Go
// cryptography and only reports its state. The FIPS build, produced with
// GOFIPS140=v1.0.0 and -tags fips140, requires the certified module and refuses
// to run without it; see Check.
package fips

import (
	"crypto/fips140"
	"errors"
	"fmt"
	"runtime/debug"
)

// CertifiedModuleVersion is the formal module version reported by
// crypto/fips140.Version for CMVP certificate #5247.
const CertifiedModuleVersion = "v1.0.0"

// CertifiedModuleSnapshot is the exact frozen snapshot GOFIPS140 records for
// that module. Matched in full rather than as a v1.0.0 prefix, so a different
// snapshot cannot be reported as the certified one.
const CertifiedModuleSnapshot = "v1.0.0-c2097c7c"

// Status describes the FIPS 140-3 state of this process.
type Status struct {
	Enabled bool

	// ModuleVersion is "latest" unless the binary was built with GOFIPS140, so
	// Enabled alone does not prove the certified module is in use.
	ModuleVersion string

	// BuildSetting is the GOFIPS140 value recorded at build time, if any.
	BuildSetting string
}

// Current inspects the running binary and returns its FIPS 140-3 status.
func Current() Status {
	s := Status{
		Enabled:       fips140.Enabled(),
		ModuleVersion: fips140.Version(),
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "GOFIPS140" {
				s.BuildSetting = setting.Value
				break
			}
		}
	}
	return s
}

// BuiltForFIPS reports whether the certified module was compiled in.
func (s Status) BuiltForFIPS() bool {
	return s.BuildSetting == CertifiedModuleSnapshot
}

// Certified reports whether the certified module is both compiled in and active.
func (s Status) Certified() bool {
	return s.Enabled && s.BuiltForFIPS()
}

// Check reports whether this build is running the cryptography it claims.
// It always returns nil in the standard build. In the FIPS build it returns an
// error unless the certified module is compiled in and active, so a binary
// carrying the fips140 tag can never serve traffic on unvalidated cryptography.
func Check() error {
	if !RequireCertified {
		return nil
	}
	s := Current()
	switch {
	case s.Certified():
		return nil
	case !s.BuiltForFIPS():
		return fmt.Errorf("built with -tags fips140 but not against the certified module "+
			"(GOFIPS140=%q, want %q); rebuild with GOFIPS140=%s",
			s.BuildSetting, CertifiedModuleSnapshot, CertifiedModuleVersion)
	default:
		return errors.New("FIPS 140-3 mode is disabled at runtime; remove fips140=off from GODEBUG")
	}
}

func (s Status) String() string {
	switch {
	case s.Certified():
		return fmt.Sprintf("enabled (Go Cryptographic Module %s, GOFIPS140=%s)", s.ModuleVersion, s.BuildSetting)
	case s.BuiltForFIPS():
		return fmt.Sprintf("disabled (built with GOFIPS140=%s but GODEBUG=fips140=off at startup)", s.BuildSetting)
	case s.Enabled:
		return fmt.Sprintf("enabled, uncertified (Go Cryptographic Module %s; binary not built with GOFIPS140)", s.ModuleVersion)
	default:
		return "disabled (binary not built with GOFIPS140)"
	}
}
