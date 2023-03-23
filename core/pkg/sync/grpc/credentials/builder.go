package credentials

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// Prefix for GRPC URL inputs. GRPC does not define a standard prefix. This prefix helps to differentiate remote
	// URLs for REST APIs (i.e - HTTP) from GRPC endpoints.
	Prefix       = "grpc://"
	PrefixSecure = "grpcs://"

	tlsVersion = tls.VersionTLS12
)

type Builder interface {
	Build(string, string) (credentials.TransportCredentials, error)
}

type CredentialBuilder struct{}

// Build is a helper to build grpc credentials.TransportCredentials based on source and cert path
func (cb *CredentialBuilder) Build(source string, certPath string) (credentials.TransportCredentials, error) {
	if strings.Contains(source, Prefix) {
		return insecure.NewCredentials(), nil
	}

	if !strings.Contains(source, PrefixSecure) {
		return nil, fmt.Errorf("invalid source. grpc source must contain prefix %s or %s", Prefix, PrefixSecure)
	}

	if certPath == "" {
		// Rely on CA certs provided from system
		return credentials.NewTLS(&tls.Config{MinVersion: tlsVersion}), nil
	}

	// Rely on provided certificate
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(certBytes) {
		return nil, fmt.Errorf("invalid certificate provided at path: %s", certPath)
	}

	return credentials.NewTLS(&tls.Config{
		MinVersion: tlsVersion,
		RootCAs:    cp,
	}), nil
}
