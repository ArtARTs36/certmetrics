package storage

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type Env struct{}

func NewEnv() *Env {
	return &Env{}
}

func (h *Env) ListFiles(_ context.Context, path string) ([]string, error) {
	return []string{path}, nil
}

func (h *Env) ReadFile(_ context.Context, path string) ([]byte, error) {
	varName := strings.TrimPrefix(path, "env://")
	varName = strings.TrimPrefix(varName, "$")

	val, ok := os.LookupEnv(varName)
	if !ok {
		return nil, fmt.Errorf("environment variable %q not found", varName)
	}

	return []byte(val), nil
}
