package integration_test

import (
	"testing"

	"github.com/open-feature/go-sdk-contrib/tests/flagd/pkg/integration"

	"github.com/cucumber/godog"
)

const flagConfigurationPath = "../../test-harness/testing-flags.json"

func TestCaching(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	initializeCachingScenario, err := integration.InitializeCachingScenario(flagConfigurationPath)
	if err != nil {
		t.Fatal(err)
	}

	suite := godog.TestSuite{
		Name:                "caching.feature",
		ScenarioInitializer: initializeCachingScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../test-harness/features/caching.feature"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run caching tests")
	}
}
