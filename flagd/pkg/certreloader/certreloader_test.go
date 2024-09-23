package certreloader

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"testing"
	"time"
)

func TestNewCertReloader(t *testing.T) {
	cert1, key1, cleanup := generateValidCertificateFiles(t)
	defer cleanup()
	_, key2, cleanup := generateValidCertificateFiles(t)
	defer cleanup()

	tcs := []struct {
		name   string
		config Config
		err    error
	}{
		{
			name:   "no config set",
			config: Config{},
			err:    fmt.Errorf("failed to load initial certificate: failed to load key pair: open : no such file or directory"),
		},
		{
			name:   "invalid certs",
			config: Config{CertPath: cert1, KeyPath: key2},
			err:    fmt.Errorf("failed to load initial certificate: failed to load key pair: tls: private key does not match public key"),
		},

		{
			name:   "valid certs",
			config: Config{CertPath: cert1, KeyPath: key1},
			err:    nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			reloader, err := NewCertReloader(tc.config)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("NewCertReloader returned error when no error was expected: %s", err)
				} else if tc.err.Error() != err.Error() {
					t.Fatalf("expected error did not matched received error. expected: %v, received: %v", tc.err, err)
				}
			} else {
				if reloader == nil {
					t.Fatal("expected reloader to not be nil")
				}
			}
		})
	}
}

func TestCertificateReload(t *testing.T) {
	newCert, newKey, cleanup := generateValidCertificateFiles(t)
	defer cleanup()

	tcs := []struct {
		name           string
		waitInterval   time.Duration
		reloadInterval time.Duration
		newCert        string
		newKey         string
		shouldRotate   bool
		err            error
	}{
		{
			name:           "reloads after interval",
			waitInterval:   time.Microsecond * 200,
			reloadInterval: time.Microsecond * 100,
			newCert:        newCert,
			newKey:         newKey,
			shouldRotate:   true,
			err:            nil,
		},
		{
			name:           "doesnt reload before  interval",
			waitInterval:   time.Microsecond * 50,
			reloadInterval: time.Microsecond * 100,
			newCert:        newCert,
			newKey:         newKey,
			shouldRotate:   false,
			err:            nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cert, key, cleanup := generateValidCertificateFiles(t)
			defer cleanup()
			reloader, err := NewCertReloader(Config{
				CertPath:       cert,
				KeyPath:        key,
				ReloadInterval: tc.reloadInterval,
			})
			if err != nil {
				t.Fatal(err)
			}

			if err := copyFile(tc.newCert, cert); err != nil {
				t.Fatalf("failed to move %s -> %s: %s", newCert, cert, err)
			}
			if err := copyFile(tc.newKey, key); err != nil {
				t.Fatalf("failed to move %s -> %s: %s", newKey, key, err)
			}
			time.Sleep(tc.waitInterval)

			actualCert, err := reloader.GetCertificate()
			if err != nil {
				t.Fatal(err)
			}
			actualCertParsed, err := x509.ParseCertificate(actualCert.Certificate[0])
			if err != nil {
				t.Fatal(err)
			}

			var expectedCert tls.Certificate
			if tc.shouldRotate {
				expectedCert, err = tls.LoadX509KeyPair(tc.newCert, tc.newKey)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				expectedCert, err = tls.LoadX509KeyPair(cert, key)
				if err != nil {
					t.Fatal(err)
				}
			}
			expectedCertParsed, err := x509.ParseCertificate(expectedCert.Certificate[0])
			if err != nil {
				t.Fatal(err)
			}
			if expectedCertParsed.DNSNames[0] != actualCertParsed.DNSNames[0] {
				t.Fatalf("expected certificate was not returned by GetCertificate. expectedCert: %v, actualCert: %v", expectedCertParsed.DNSNames[0], actualCertParsed.DNSNames[0])
			}
		})
	}
}

func generateValidCertificate(t *testing.T) (*bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	// set up our CA certificate
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatal(err)
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		t.Fatal(err)
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		t.Fatal(err)
	}

	caPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	if err != nil {
		t.Fatal(err)
	}

	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		DNSNames:     []string{randString(8)},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatalf("failed to create private key: %s", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		t.Fatalf("failed to create certificate: %s", err)
	}

	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		t.Fatal(err)
	}

	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		t.Fatal(err)
	}

	return certPEM, certPrivKeyPEM
}

func generateValidCertificateFiles(t *testing.T) (string, string, func()) {
	t.Helper()
	certFile, err := os.CreateTemp("", "certreloader_cert")
	if err != nil {
		t.Fatalf("failed to create certFile: %s", err)
	}
	defer certFile.Close()
	keyFile, err := os.CreateTemp("", "certreloader_key")
	if err != nil {
		t.Fatalf("failed to create keyFile: %s", err)
	}
	defer keyFile.Close()

	certBytes, keyBytes := generateValidCertificate(t)
	if _, err := io.Copy(certFile, certBytes); err != nil {
		t.Fatalf("failed to copy certBytes into %s: %s", certFile.Name(), err)
	}
	if _, err := io.Copy(keyFile, keyBytes); err != nil {
		t.Fatalf("failed to copy keyBytes into %s: %s", keyFile.Name(), err)
	}

	return certFile.Name(), keyFile.Name(), func() {
		os.Remove(certFile.Name())
		os.Remove(keyFile.Name())
	}
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to load key pair: %w", err)
	}

	err = os.WriteFile(dst, data, 0o0600)
	if err != nil {
		return fmt.Errorf("failed to load key pair: %w", err)
	}
	return nil
}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, n)
	//nolint:errcheck
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
