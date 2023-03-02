package integration_test

import (
	"flag"
	"testing"

	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"

	"github.com/cucumber/godog"
	"github.com/open-feature/go-sdk-contrib/tests/flagd/pkg/integration"
)

func TestEvaluation(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	flag.Parse()

	var providerOptions []flagd.ProviderOption
	name := "evaluation.feature"

	if tls == "true" {
		name = "evaluation_tls.feature"
		providerOptions = []flagd.ProviderOption{flagd.WithTLS(certPath)}
	}

	testSuite := godog.TestSuite{
		Name:                name,
		ScenarioInitializer: integration.InitializeEvaluationScenario(providerOptions...),
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../../test-harness/features/evaluation.feature"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if testSuite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run evaluation tests")
	}
}
