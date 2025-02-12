package config

import (
	"testing"

	"github.com/matryer/is"
	"github.com/redis/go-redis/v9"
)

func TestRedisOptions(t *testing.T) {
	is := is.New(t)

	expected := &redis.Options{
		Network:  "tcp",
		Addr:     "localhost:6379",
		Password: "password",
		DB:       1,
	}

	r := Redis{
		URI: "redis://:password@localhost/1",
	}

	opts, err := r.Options()

	is.NoErr(err)
	is.Equal(expected, opts)
}
