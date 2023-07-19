package main

import (
	"fmt"
	"os"
	"strconv"

	pb "buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	"google.golang.org/grpc"
)

func main() {

	waitSecondsBetweenRequests := getWaitSecondsBetweenRequests()
	flagDTURL := getURL()

	// Create a channel to receive a signal when the gRPC connection fails
	errChan := make(chan bool)

	// Use a goroutine to run your program logic
	go func() {
		if err := handleRequests(waitSecondsBetweenRequests, flagDTURL); err != nil {
			errChan <- true
		}
	}()

	// Wait for the specified duration or until the program completes
	<-errChan
}

func handleRequests(waitSecondsBetweenRequests int, flagDTURL string) error {
	var conn *grpc.ClientConn
	var grpcClient pb.FlagSyncServiceClient
	// open the connection only once
	conn, grpcClient = establishGrpcConnection(flagDTURL)

	defer func() {
		if conn != nil {
			// clean up
			err := conn.Close()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()

	return doRequests(grpcClient, waitSecondsBetweenRequests)
}

func getWaitSecondsBetweenRequests() int {
	return getEnvVarOrDefault("WAIT_TIME_BETWEEN_REQUESTS_MS", 1)
}

func getURL() string {
	return getEnvOrDefault("URL", "flagd-proxy-svc.flagdt-dev:8015")
}

func getEnvVarOrDefault(envVar string, defaultValue int) int {
	if envVarValue := os.Getenv(envVar); envVarValue != "" {
		parsedEnvVarValue, err := strconv.ParseInt(envVarValue, 10, 64)
		if err == nil && parsedEnvVarValue > 0 {
			defaultValue = int(parsedEnvVarValue)
		}
	}
	return defaultValue
}

func getEnvOrDefault(envVar string, defaultValue string) string {
	if envVarValue := os.Getenv(envVar); envVarValue != "" {
		return envVarValue
	}
	return defaultValue
}
