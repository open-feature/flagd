package service

import (
	"crypto/rand"
	"crypto/tls"

	"google.golang.org/grpc/credentials"
)

func loadTLSCredentials(serverCertPath string, serverKeyPath string) (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	creds, err := credentials.NewServerTLSFromFile(serverCertPath, serverKeyPath)
	if err != nil {
		return nil, err
	}

	return creds, nil
}

func loadTLSConfig(certPath, keyPath string) (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair(certPath,
		keyPath)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		Rand:         rand.Reader,
		MinVersion:   tls.VersionTLS12,
	}

	return config, nil
}
