package service

import (
	"crypto/rand"
	"crypto/tls"
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

func loadTLSCertificate(certPath, keyPath string) (*tls.Certificate, error) {
	certificate, err := tls.LoadX509KeyPair(certPath,
		keyPath)
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		Rand:         rand.Reader,
		MinVersion:   tls.VersionTLS12,
	}

	return &config.Certificates[0], nil
}
