package storage

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/store_mock.go -package=mocks github.com/Jsee98/Recallr/storage Store
type Store interface {
	// Basic key-value operations
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error

	HSet(ctx context.Context, key string, field string, value string) error
	HGet(ctx context.Context, key string, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, field string) error

	LPush(ctx context.Context, key string, value string) error
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)

	Exists(ctx context.Context, key string) (bool, error)
	Type() string

	Close() error
}
