package storage

import (
	"context"
)

type resolvable struct {
	resolver *Resolver
}

func Resolve(resolver *Resolver) Storage {
	return &resolvable{
		resolver: resolver,
	}
}

func (r *resolvable) ListFiles(ctx context.Context, path string) ([]string, error) {
	return r.resolver.Resolve(path).ListFiles(ctx, path)
}

func (r *resolvable) ReadFile(ctx context.Context, path string) ([]byte, error) {
	return r.resolver.Resolve(path).ReadFile(ctx, path)
}
