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

var testSource1 = "testSource1"
var testSource2 = "testSource2"
var testVariants = map[string]any{
	"true":  true,
	"false": false,
}
var testSource1Flags = []model.Flag{
	{
		Key:            "flagA",
		State:          "ENABLED",
		DefaultVariant: "false",
		Variants:       testVariants,
	},
}
var testSource2Flags = []model.Flag{
	{
		Key:            "flagB",
		State:          "ENABLED",
		DefaultVariant: "true",
		Variants:       testVariants,
	},
}

// getSimpleFlagStore is a test util which returns a flag store pre-filled with flags from sources testSource1 and testSource2.
func getSimpleFlagStore(t testing.TB) (store.IStore, []string) {
	t.Helper()

	sources := []string{testSource1, testSource2}

	flagStore, err := store.NewStore(logger.NewLogger(nil, false), sources)
	if err != nil {
		t.Fatalf("error creating flag store: %v", err)
	}

	flagStore.Update(testSource1, testSource1Flags, model.Metadata{
		"keyDuped": "value",
		"keyA":     "valueA",
	})

	flagStore.Update(testSource2, testSource2Flags, model.Metadata{
		"keyDuped": "value",
		"keyB":     "valueB",
	})

	return flagStore, sources
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
