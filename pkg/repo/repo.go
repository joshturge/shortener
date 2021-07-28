package repo

import (
	"context"
	"errors"
	"io"
)

var ErrNotFound = errors.New("value was not found in repository")

// Repository is a key value datastore where urls and there corresponding hashs are kept
type Repository interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
	io.Closer
}
