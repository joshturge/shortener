package repo

import (
	"context"
	"time"
)

type url struct {
	str string
	exp int64
}

var repo map[string]url

func init() {
	repo = make(map[string]url)
}

// Map is a in-memory keystore which implements the Repository interface which is used for testing
type Map struct {
	timeout time.Duration
}

// NewMap will create a new repository
func NewMap() Repository {
	return &Map{}
}

// SetTimeout will set the timeout for keys
func (m *Map) SetTimeout(d time.Duration) {
	m.timeout = d
}

// Get a value out of a map
func (m *Map) Get(ctx context.Context, key string) (string, error) {
	if val, ok := repo[key]; ok {
		if time.Now().Unix() < val.exp {
			return val.str, nil
		} else {
			delete(repo, key)
		}
	}

	return "", ErrNotFound
}

// Set a value within the map
func (m *Map) Set(ctx context.Context, key, val string) error {
	repo[key] = url{str: val, exp: time.Now().Add(m.timeout).Unix()}
	return nil
}

// Close dummy, GC will clean this up
func (m Map) Close() error {
	return nil
}
