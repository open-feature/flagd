package credentials

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const tlsVersion = tls.VersionTLS12

type Builder interface {
	Build(secure bool, certPath string) (credentials.TransportCredentials, error)
}

type CredentialBuilder struct{}

// Build is a helper to build grpc credentials.TransportCredentials based on source and cert path
func (cb *CredentialBuilder) Build(secure bool, certPath string) (credentials.TransportCredentials, error) {
	if !secure {
		return insecure.NewCredentials(), nil
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
