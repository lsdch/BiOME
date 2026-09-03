package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type RawFileStorage interface {
	Store(ctx context.Context, key string, r io.Reader) (checksum string, err error)
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}

type FilesystemRawFileStorage struct {
	root string
}

func NewFilesystemRawFileStorage(root string) (*FilesystemRawFileStorage, error) {
	if root == "" {
		return nil, fmt.Errorf("raw file storage root cannot be empty")
	}

	if err := os.MkdirAll(root, 0o750); err != nil {
		return nil, fmt.Errorf("create raw file storage root: %w", err)
	}

	return &FilesystemRawFileStorage{
		root: root,
	}, nil
}
func (s *FilesystemRawFileStorage) Store(ctx context.Context, key string, r io.Reader) (string, error) {
	path, err := s.path(key)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return "", fmt.Errorf("create raw file directory: %w", err)
	}

	logrus.Debugf("Storing raw file at %s", path)

	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0o640,
	)
	if err != nil {
		return "", fmt.Errorf("create raw file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()

	writer := io.MultiWriter(file, hash)

	if _, err := io.Copy(writer, r); err != nil {
		return "", fmt.Errorf("write raw file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (s *FilesystemRawFileStorage) Open(
	ctx context.Context,
	key string,
) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	path, err := s.path(key)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open raw file: %w", err)
	}

	return file, nil
}

func (s *FilesystemRawFileStorage) Delete(
	ctx context.Context,
	key string,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	path, err := s.path(key)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("delete raw file: %w", err)
	}

	return nil
}

func (s *FilesystemRawFileStorage) path(key string) (string, error) {
	key = filepath.ToSlash(key)
	key = strings.TrimPrefix(key, "/")

	if key == "" || key == "." {
		return "", fmt.Errorf("invalid raw file storage key")
	}

	// Prevent escaping the storage root with "../".
	clean := filepath.Clean(filepath.FromSlash(key))
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("invalid raw file storage key: %q", key)
	}

	return filepath.Join(s.root, clean), nil
}
