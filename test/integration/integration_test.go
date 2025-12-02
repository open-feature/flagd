package integration_test

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/open-feature/go-sdk-contrib/tests/flagd/testframework"
	"github.com/testcontainers/testcontainers-go"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func getGitCommitNative() string {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "no-git-repo"
	}

	ref, err := repo.Head()
	if err != nil {
		return "no-head"
	}

	hash := ref.Hash().String()
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}

func getPlatform() string {
	switch runtime.GOARCH {
	case "amd64":
		return "linux/amd64"
	case "arm64":
		return "linux/arm64"
	default:
		return "linux/amd64" // Always return a valid platform
	}
}

func TestMain(m *testing.M) {
	flag.Parse()

	// Setup: Build and start flagd-testbed container
	ctx := context.Background()
	var err error

	err = setupFlagdTestbed(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to setup flagd testbed: %v", err))
	}

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func setupFlagdTestbed(ctx context.Context) error {
	// Build the testbed image with local flagd
	buildContext, err := filepath.Abs("../../") // Adjust path to your project root
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Build flagd first
	flagdImage, err := buildFlagdImage(ctx, buildContext)
	if err != nil {
		return fmt.Errorf("failed to build flagd image: %w", err)
	}

	// Build testbed image using the local flagd
	err = buildTestbedImage(ctx, buildContext, flagdImage)
	if err != nil {
		return fmt.Errorf("failed to build testbed image: %w", err)
	}

	return nil
}

func buildFlagdImage(ctx context.Context, buildContext string) (string, error) {
	imageName := "flagd:test-local"

	// Execute the docker buildx command
	cmd := exec.Command("docker", "buildx", "build",
		"--platform="+getPlatform(),
		"-t", imageName,
		"-f", "flagd/build.Dockerfile",
		".")

	cmd.Dir = buildContext

	// Print the exact command being run
	fmt.Printf("Running command: %s\n", strings.Join(cmd.Args, " "))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return imageName, fmt.Errorf("docker buildx failed: %v\nOutput: %s", err, output)
	}

	return imageName, nil
}

func buildTestbedImage(ctx context.Context, buildContext, flagdImage string) error {
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    filepath.Join(buildContext, "./test-harness"), // Path to testbed directory
			Dockerfile: "./flagd/Dockerfile",
			BuildArgs: map[string]*string{
				"FLAGD_BASE_IMAGE": &flagdImage,
			},
			Tag:  getGitCommitNative(),
			Repo: "local-testbed-image",
		},
	}

	// Build the image
	_, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false, // We just want to build, not start
	})
	if err != nil {
		return fmt.Errorf("failed to build testbed image: %w", err)
	}

	return nil
}

func TestRPC(t *testing.T) {

	// Setup testbed runner for RPC provider
	runner := testframework.NewTestbedRunner(testframework.TestbedConfig{
		ResolverType:  testframework.RPC,
		TestbedConfig: "default", // Use default testbed configuration
		TestbedDir:    "../../test-harness",
		Image:         "local-testbed-image",
		Tag:           getGitCommitNative(),
	})
	defer runner.Cleanup()

	// Define feature paths - using flagd-testbed gherkin files
	featurePaths := []string{
		"./",
	}

	// Run tests with RPC-specific tags - exclude connection/event issues we won't tackle
	tags := "@rpc && ~@unixsocket && ~@targetURI && ~@sync && ~@metadata && ~@grace && ~@events && ~@customCert && ~@reconnect && ~@caching && ~@forbidden"

	if err := runner.RunGherkinTestsWithSubtests(t, featurePaths, tags); err != nil {
		t.Fatalf("Gherkin tests failed: %v", err)
	}
}

func TestInProcess(t *testing.T) {
	// Setup testbed runner for RPC provider
	runner := testframework.NewTestbedRunner(testframework.TestbedConfig{
		ResolverType:  testframework.InProcess,
		TestbedConfig: "default", // Use default testbed configuration
		TestbedDir:    "../../test-harness",
		Image:         "local-testbed-image",
		Tag:           getGitCommitNative(),
	})
	defer runner.Cleanup()

	// Define feature paths - using flagd-testbed gherkin files
	featurePaths := []string{
		"./",
	}

	// Run tests with InProcess-specific tags
	tags := "@in-process && ~@unixsocket&& ~@metadata && ~@contextEnrichment && ~@customCert && ~@forbidden && ~@sync-port && ~@sync-payload"

	if err := runner.RunGherkinTestsWithSubtests(t, featurePaths, tags); err != nil {
		t.Fatalf("Gherkin tests failed: %v", err)
	}
}
