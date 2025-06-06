package dragonfly

import (
	"context"
	"time"

	"github.com/JSee98/Recallr/constants"
	"github.com/redis/go-redis/v9"
)

type DragonflyStore struct {
	client *redis.Client
}

func NewDragonflyStore(config *DragonflyConfig) *DragonflyStore {

	redisConfig := config.toRedisOptions()

	rdb := redis.NewClient(redisConfig)
	return &DragonflyStore{client: rdb}
}

// Implement Store interface
func (ds *DragonflyStore) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return ds.client.Set(ctx, key, value, expiration).Err()
}

func (ds *DragonflyStore) Get(ctx context.Context, key string) (string, error) {
	return ds.client.Get(ctx, key).Result()
}

func (ds *DragonflyStore) Delete(ctx context.Context, key string) error {
	return ds.client.Del(ctx, key).Err()
}

func (ds *DragonflyStore) HSet(ctx context.Context, key string, field string, value string) error {
	return ds.client.HSet(ctx, key, field, value).Err()
}

func (ds *DragonflyStore) HGet(ctx context.Context, key string, field string) (string, error) {
	return ds.client.HGet(ctx, key, field).Result()
}

func (ds *DragonflyStore) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return ds.client.HGetAll(ctx, key).Result()
}

func (ds *DragonflyStore) HDel(ctx context.Context, key string, field string) error {
	return ds.client.HDel(ctx, key, field).Err()
}

func (ds *DragonflyStore) LPush(ctx context.Context, key string, value string) error {
	return ds.client.LPush(ctx, key, value).Err()
}

func (ds *DragonflyStore) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return ds.client.LRange(ctx, key, start, stop).Result()
}

func (ds *DragonflyStore) Exists(ctx context.Context, key string) (bool, error) {
	res, err := ds.client.Exists(ctx, key).Result()
	return res > 0, err
}

func (ds *DragonflyStore) Close() error {
	return ds.client.Close()
}

func (ds *DragonflyStore) Type() string {
	return constants.DragonFlyType
}
