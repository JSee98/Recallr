package dragonfly

import (
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
)

type DragonflyConfig struct {
	Addr     string
	Username string
	Password string
	DB       int

	MaxRetries      *int
	DialTimeout     *time.Duration
	ReadTimeout     *time.Duration
	WriteTimeout    *time.Duration
	PoolSize        *int
	MinIdleConns    *int
	TLSConfig       *tls.Config
	ClientName      *string
	DisableIdentity *bool
}

func (cfg *DragonflyConfig) toRedisOptions() *redis.Options {
	opts := &redis.Options{
		Addr:      cfg.Addr,
		Username:  cfg.Username,
		Password:  cfg.Password,
		DB:        cfg.DB,
		TLSConfig: cfg.TLSConfig,
	}

	if cfg.MaxRetries != nil {
		opts.MaxRetries = *cfg.MaxRetries
	}
	if cfg.DialTimeout != nil {
		opts.DialTimeout = *cfg.DialTimeout
	}
	if cfg.ReadTimeout != nil {
		opts.ReadTimeout = *cfg.ReadTimeout
	}
	if cfg.WriteTimeout != nil {
		opts.WriteTimeout = *cfg.WriteTimeout
	}
	if cfg.PoolSize != nil {
		opts.PoolSize = *cfg.PoolSize
	}
	if cfg.MinIdleConns != nil {
		opts.MinIdleConns = *cfg.MinIdleConns
	}
	if cfg.ClientName != nil {
		opts.ClientName = *cfg.ClientName
	}
	if cfg.DisableIdentity != nil {
		opts.DisableIdentity = *cfg.DisableIdentity
	}

	return opts
}
