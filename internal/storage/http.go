package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HTTP struct {
}

func NewHTTP() *HTTP {
	return &HTTP{}
}

func (h *HTTP) ListFiles(_ context.Context, path string) ([]string, error) {
	return []string{path}, nil
}

func (h *HTTP) ReadFile(ctx context.Context, path string) ([]byte, error) {
	uri, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	req := &http.Request{
		URL: uri,
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return bs, nil
}
