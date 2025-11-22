package storage

import "context"

type Storage interface {
	ListFiles(ctx context.Context, dir string) ([]string, error)
	ReadFile(ctx context.Context, path string) ([]byte, error)
}
