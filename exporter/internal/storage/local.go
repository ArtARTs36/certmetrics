package storage

import (
	"context"
	"os"
	"path/filepath"
)

type Local struct{}

func NewLocal() *Local {
	return &Local{}
}

func (l *Local) ListFiles(_ context.Context, dir string) ([]string, error) {
	return filepath.Glob(dir)
}

func (l *Local) ReadFile(_ context.Context, path string) ([]byte, error) {
	return os.ReadFile(path)
}
