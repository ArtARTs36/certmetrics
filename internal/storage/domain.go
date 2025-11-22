package storage

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
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

func (h *Domain) ListFiles(ctx context.Context, path string) ([]string, error) {
	uri, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse path %q: %w", path, err)
	}

	port := uri.Port()
	if port == "" {
		port = "443"
	}

	newPath := fmt.Sprintf("%s:%s", uri.Host, port)

	conn, err := h.dial(ctx, newPath)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			slog.WarnContext(ctx, "error closing connection", slog.String("addr", newPath), slog.Any("err", err))
		}
	}()

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

func (h *Domain) dial(ctx context.Context, addr string) (*tls.Conn, error) {
	dialer := tls.Dialer{
		Config: &tls.Config{
			InsecureSkipVerify: true, //nolint: gosec // not need
		},
	}

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("tcp dial: %w", err)
	}

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		slog.WarnContext(ctx, "error closing connection", slog.String("addr", addr), slog.Any("err", err))
		return nil, errors.New("tcp conn does not implement tls.Conn")
	}

	return tlsConn, nil
}
