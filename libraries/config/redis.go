package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis configures a connection to a Redis database.
type Redis struct {
	URI          string        `env:"URI"`
	DialTimeout  time.Duration `env:"DIAL_TIMEOUT"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT"`
}

// Options returns a configured redis.Options structure.
func (r Redis) Options() (*redis.Options, error) {
	opts, err := redis.ParseURL(r.URI)
	if err != nil {
		return nil, err
	}

	opts.DialTimeout = r.DialTimeout
	opts.ReadTimeout = r.ReadTimeout
	opts.WriteTimeout = r.WriteTimeout

	return opts, nil
}

// Connect returns a connected redis.Client instance.
func (r Redis) Connect(ctx context.Context) (*redis.Client, error) {
	opts, err := r.Options()
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return client, err
	}

	return client, nil
}
