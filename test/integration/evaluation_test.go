package integration_test

import (
	"flag"
	"testing"

	"github.com/cucumber/godog"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk-contrib/tests/flagd/pkg/integration"
	"github.com/open-feature/go-sdk/pkg/openfeature"
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
		Name: name,
		ScenarioInitializer: integration.InitializeEvaluationScenario(func() openfeature.FeatureProvider {
			return flagd.NewProvider(providerOptions...)
		}),
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../spec/specification/assets/gherkin/evaluation.feature"},
			TestingT: t, // Testing instance that will run subtests.
			Strict:   true,
		},
	}

	if testSuite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run evaluation tests")
	}
}
