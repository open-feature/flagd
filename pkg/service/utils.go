package service

import (
	"log"

	"google.golang.org/grpc/credentials"
)

func loadTLSCredentials(serverCertPath string, serverKeyPath string) (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	creds, err := credentials.NewServerTLSFromFile(serverCertPath, serverKeyPath)
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}

	return creds, nil
}
