package main

import (
	"fmt"
	"os"
	"strconv"

	pb "buf.build/gen/go/open-feature-forking/flagd/grpc/go/sync/v1/syncv1grpc"
	"google.golang.org/grpc"
)

func main() {
	waitSecondsBetweenRequests := getWaitSecondsBetweenRequests()
	flagdURL := getURL()

	// Create a channel to receive a signal when the gRPC connection fails
	errChan := make(chan bool)

	// Use a goroutine to run your program logic
	go func() {
		if err := handleRequests(waitSecondsBetweenRequests, flagdURL); err != nil {
			errChan <- true
		}
	}()

	// The program should run until it receives an error
	<-errChan
}

func handleRequests(waitSecondsBetweenRequests int, flagdURL string) error {
	var conn *grpc.ClientConn
	var grpcClient pb.FlagSyncServiceClient
	// open the connection only once
	conn, grpcClient = establishGrpcConnection(flagdURL)

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
	return getEnvVarOrDefault("WAIT_TIME_BETWEEN_REQUESTS_S", 1)
}

func getURL() string {
	return getEnvOrDefault("URL", "flagd-proxy-svc.flagd-dev:8015")
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
