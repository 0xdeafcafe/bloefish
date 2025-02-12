package config

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	URI          string `env:"URI"`
	DatabaseName string `env:"DATABASE_NAME"`
}

func (m MongoDB) MustConnect(ctx context.Context) (*mongo.Client, *mongo.Database) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.URI))
	if err != nil {
		panic(fmt.Errorf("failed to connect to mongodb: %w", err))
	}

	if err := client.Ping(ctx, nil); err != nil {
		panic(fmt.Errorf("failed to ping mongodb: %w", err))
	}

	return client, client.Database(m.DatabaseName)
}
