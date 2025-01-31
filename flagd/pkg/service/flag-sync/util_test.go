package sync

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"google.golang.org/grpc/credentials"
)

// getSimpleFlagStore returns a flag store pre-filled with flags from sources A & B & C, which C empty
func getSimpleFlagStore() (*store.State, []string) {
	variants := map[string]any{
		"true":  true,
		"false": false,
	}

	flagStore := store.NewFlags()

	flagStore.Set("flagA", model.Flag{
		State:          "ENABLED",
		DefaultVariant: "false",
		Variants:       variants,
		Source:         "A",
	})

	flagStore.Set("flagB", model.Flag{
		State:          "ENABLED",
		DefaultVariant: "true",
		Variants:       variants,
		Source:         "B",
	})

	return flagStore, []string{"A", "B", "C"}
}

func loadTLSClientCredentials(certPath string) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from path '%s'", certPath)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}

	return credentials.NewTLS(config), nil
}
