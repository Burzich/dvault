package fs

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/Burzich/dvault/internal/dvault/storage"
)

type Storage struct {
	mountPoint string
}

func NewFSStorage(mountPath string) *Storage {
	return &Storage{
		mountPoint: mountPath,
	}
}

func (f Storage) Put(_ context.Context, path string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(filepath.Join(f.mountPoint, path)), os.ModePerm)
	if err != nil {
		return err
	}

	p := filepath.Join(f.mountPoint, path)

	err = os.WriteFile(p, data, 0644)
	if errors.Is(err, os.ErrNotExist) {
		return storage.ErrPathNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

func (f Storage) Get(_ context.Context, path string) ([]byte, error) {
	p := filepath.Join(f.mountPoint, path)

	data, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return nil, storage.ErrPathNotFound
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (f Storage) Delete(_ context.Context, path string) error {
	p := filepath.Join(f.mountPoint, path)

	err := os.Remove(p)
	if errors.Is(err, os.ErrNotExist) {
		return storage.ErrPathNotFound
	}
	if err != nil {
		return err
	}

	return nil
}
