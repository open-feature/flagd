package sync

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"google.golang.org/grpc/credentials"
)

// getSimpleFlagStore is a test util which returns a flag store pre-filled with flags from sources A & B & C, which C empty
func getSimpleFlagStore(t *testing.T) (*store.Store, []string) {
	variants := map[string]any{
		"true":  true,
		"false": false,
	}

	flagStore, err := store.NewStore(logger.NewLogger(nil, false))
	if err != nil {
		t.Fatalf("error creating flag store: %v", err)
	}

	flagStore.Update("A", "", map[string]model.Flag{
		"flagA": {
			State:          "ENABLED",
			DefaultVariant: "false",
			Variants:       variants,
		},
	}, model.Metadata{
		"keyDuped": "value",
		"keyA":     "valueA",
	})

	flagStore.Update("B", "", map[string]model.Flag{
		"flagB": {
			State:          "ENABLED",
			DefaultVariant: "true",
			Variants:       variants,
		},
	}, model.Metadata{
		"keyDuped": "value",
		"keyB":     "valueB",
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
