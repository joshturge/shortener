package repo

import (
	"context"
	"time"
)

type url struct {
	str string
	exp int64
}

// Map is a in-memory keystore which implements the Repository interface which is used for testing
type Map struct {
	Timeout time.Duration
	repo map[string]url
}

// NewMap will create a new repository
func NewMap(timeout time.Duration) Repository {
	return &Map{Timeout: timeout, repo: make(map[string]url)}
}

// Get a value out of a map
func (m *Map) Get(ctx context.Context, key string) (string, error) {
	if val, ok := m.repo[key]; ok {
		if time.Now().Unix() < val.exp {
			return val.str, nil
		} else {
			delete(m.repo, key)
		}
	}

	return "", ErrNotFound
}

// Set a value within the map
func (m *Map) Set(ctx context.Context, key, val string) error {
	m.repo[key] = url{str: val, exp: time.Now().Add(m.Timeout).Unix()}
	return nil
}

func (m Map) Close() error {
	return nil
}
