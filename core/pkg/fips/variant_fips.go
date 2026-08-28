//go:build fips140

package fips

// Variant names this build for logs and `flagd version`.
const Variant = "fips"

// RequireCertified makes Check fail unless the certified module is active.
const RequireCertified = true
