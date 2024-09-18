package certreloader

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"
)

type Config struct {
	KeyPath        string
	CertPath       string
	ReloadInterval time.Duration
}

type certReloader struct {
	cert       *tls.Certificate
	mu         sync.RWMutex
	nextReload time.Time
	Config
}

func NewCertReloader(config Config) (*certReloader, error) {
	reloader := certReloader{
		Config: config,
	}

	reloader.mu.Lock()
	defer reloader.mu.Unlock()
	cert, err := reloader.loadCertificate()
	if err != nil {
		return nil, fmt.Errorf("failed to load initial certificate: %w", err)
	}
	reloader.cert = &cert

	return &reloader, nil
}

func (r *certReloader) GetCertificate() (*tls.Certificate, error) {
	now := time.Now()
	// Read locking here before we do the time comparison
	// If a reload is in progress this will block and we will skip reloading in the current
	// call once we can continue
	r.mu.RLock()
	shouldReload := r.ReloadInterval != 0 && r.nextReload.Before(now)
	r.mu.RUnlock()
	if shouldReload {
		// Need to release the read lock, otherwise we deadlock
		r.mu.Lock()
		defer r.mu.Unlock()
		cert, err := r.loadCertificate()
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS cert and key: %w", err)
		}
		r.cert = &cert
		r.nextReload = now.Add(r.ReloadInterval)
		return r.cert, nil
	}
	return r.cert, nil
}

func (c *certReloader) loadCertificate() (tls.Certificate, error) {
	newCert, err := tls.LoadX509KeyPair(c.CertPath, c.KeyPath)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to load key pair: %w", err)
	}

	return newCert, nil
}
