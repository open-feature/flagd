// Package fips reports whether this binary is running the CMVP-certified Go
// Cryptographic Module in FIPS 140-3 approved mode.
package fips

import (
	"crypto/fips140"
	"fmt"
	"runtime/debug"
	"strings"
)

// CertifiedModuleVersion is the module flagd builds against: CMVP certificate
// #5247. GOFIPS140 records a resolved snapshot ("v1.0.0-c2097c7c"), so this is
// matched as a prefix.
const CertifiedModuleVersion = "v1.0.0"

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

// BuiltForFIPS reports whether the certified module was compiled in. Release
// artifacts always are; a plain `go build` is not.
func (s Status) BuiltForFIPS() bool {
	return strings.HasPrefix(s.BuildSetting, CertifiedModuleVersion)
}

// Certified reports whether the certified module is both compiled in and active.
func (s Status) Certified() bool {
	return s.Enabled && s.BuiltForFIPS()
}

// Degraded reports a binary built for FIPS that is running outside FIPS mode,
// which happens when GODEBUG=fips140=off is set at process start.
func (s Status) Degraded() bool {
	return s.BuiltForFIPS() && !s.Enabled
}

func (s Status) String() string {
	switch {
	case s.Certified():
		return fmt.Sprintf("enabled (Go Cryptographic Module %s, GOFIPS140=%s)", s.ModuleVersion, s.BuildSetting)
	case s.Degraded():
		return fmt.Sprintf("disabled (built with GOFIPS140=%s but GODEBUG=fips140=off at startup)", s.BuildSetting)
	case s.Enabled:
		return fmt.Sprintf("enabled, uncertified (Go Cryptographic Module %s; binary not built with GOFIPS140)", s.ModuleVersion)
	default:
		return "disabled (binary not built with GOFIPS140)"
	}
}
