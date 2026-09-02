//go:build !fips140

package fips

// Variant names this build for logs and `flagd version`.
const Variant = "standard"

// RequireCertified is false: the standard build reports FIPS state but does not
// require it.
const RequireCertified = false
