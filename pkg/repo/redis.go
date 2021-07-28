package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis is a key value store which satisfies the Repository interface
type Redis struct {
	rdb     *redis.Client
	Timeout time.Duration
}

// NewRedisRepo will create a new Repository given the address, password of the redis instance
func NewRedisRepo(addr, pwd string, timeout time.Duration) (Repository, error) {
	repo := Redis{Timeout: timeout}
	repo.rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0,
	})

	// check our connection to the redis instance
	if _, err := repo.rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("unable to ping redis instance: %w", err)
	}

	return &repo, nil
}

// Set a key and value within the redis instance
func (r *Redis) Set(ctx context.Context, key, val string) (err error) {
	fmt.Printf("SET: %s = %s\n", key, val)
	if err = r.rdb.Set(ctx, key, val, r.Timeout).Err(); err != nil {
		return fmt.Errorf("unable to set key %s in redis repository: %w", key, err)
	}

	return nil
}

// Get a value out of the redis instance given the corresponding key
func (r *Redis) Get(ctx context.Context, key string) (val string, err error) {
	val, err = r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrNotFound
		}

		return "", fmt.Errorf("unable to get key %s in redis repository: %w", key, err)
	}

	fmt.Printf("GET: %s = %s\n", key, val)

	return val, nil
}

// Close the connection to the redis instance
func (r *Redis) Close() error {
	return r.rdb.Close()
}
