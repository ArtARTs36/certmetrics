package storage

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
)

type Domain struct {
	certs map[string]*x509.Certificate
}

func NewDomain() *Domain {
	return &Domain{
		certs: map[string]*x509.Certificate{},
	}
}

func (h *Domain) ListFiles(_ context.Context, path string) ([]string, error) {
	uri, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse path %q: %w", path, err)
	}

	port := uri.Port()
	if port == "" {
		port = "443"
	}

	newPath := fmt.Sprintf("%s:%s", uri.Host, port)

	conn, err := tls.Dial("tcp", newPath, &tls.Config{
		InsecureSkipVerify: true, //nolint: gosec // not need
	})
	if err != nil {
		return nil, fmt.Errorf("tcp dial: %w", err)
	}

	certs := conn.ConnectionState().PeerCertificates

	paths := make([]string, len(certs))
	for i, cert := range certs {
		indexPath := fmt.Sprintf("%s#%d", path, i)
		paths[i] = indexPath

		h.certs[indexPath] = cert
	}

	return paths, nil
}

func (h *Domain) ReadFile(_ context.Context, path string) ([]byte, error) {
	cert, ok := h.certs[path]
	if !ok {
		return nil, fmt.Errorf("certificate with path %q not found", path)
	}

	return cert.Raw, nil
}
