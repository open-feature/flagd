package integration_test

import (
	"flag"
	"testing"

	"github.com/cucumber/godog"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk-contrib/tests/flagd/pkg/integration"
	"github.com/open-feature/go-sdk/pkg/openfeature"
)

func TestJsonEvaluator(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	flag.Parse()

	var providerOptions []flagd.ProviderOption
	name := "flagd-json-evaluator.feature"

	testSuite := godog.TestSuite{
		Name: name,
		ScenarioInitializer: integration.InitializeFlagdJsonScenario(func() openfeature.FeatureProvider {
			return flagd.NewProvider(providerOptions...)
		}),
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../test-harness/gherkin/flagd-json-evaluator.feature"},
			TestingT: t, // Testing instance that will run subtests.
			Strict:   true,
		},
	}

	if testSuite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run evaluation tests")
	}
}
