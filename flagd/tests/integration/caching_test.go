package integration_test

import (
	"flag"
	"testing"

	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk-contrib/tests/flagd/pkg/integration"

	"github.com/cucumber/godog"
)

const (
	flagConfigurationPath           = "../../../test-harness/symlink_testing-flags.json"
	flagConfigurationTargetFilePath = "../../../test-harness/testing-flags.json"
	flagConfigurationMutatedPath    = "../../../test-harness/mutated-testing-flags.json"
)

func TestCaching(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	flag.Parse()

	var providerOptions []flagd.ProviderOption
	name := "caching.feature"

	if tls == "true" {
		name = "caching_tls.feature"
		providerOptions = []flagd.ProviderOption{flagd.WithTLS(certPath)}
	}

	initializeCachingScenario, err := integration.InitializeCachingScenario(
		flagConfigurationTargetFilePath, flagConfigurationPath, flagConfigurationMutatedPath, providerOptions...)
	if err != nil {
		t.Fatal(err)
	}

	testSuite := godog.TestSuite{
		Name:                name,
		ScenarioInitializer: initializeCachingScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../../test-harness/features/caching.feature"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if testSuite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run caching tests")
	}
}
