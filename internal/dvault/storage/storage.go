package storage

import (
	"context"
	"io"
)

type Storage interface {
	Put(ctx context.Context, path string, data []byte) error
	Get(ctx context.Context, path string) ([]byte, error)
	Delete(ctx context.Context, path string) error
	List(ctx context.Context, path string) ([]string, error)
	io.Closer
}
