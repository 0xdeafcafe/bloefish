package internal

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/libraries/telemetry"
	"github.com/0xdeafcafe/bloefish/services/fileupload/internal/app"
	"github.com/0xdeafcafe/bloefish/services/fileupload/internal/app/repositories"
	"github.com/0xdeafcafe/bloefish/services/fileupload/internal/app/services"
	"github.com/0xdeafcafe/bloefish/services/fileupload/internal/transport/rpc"
)

type Config struct {
	Server    config.Server    `env:"SERVER"`
	Telemetry telemetry.Config `env:"TELEMETRY"`
	Logging   clog.Config      `env:"LOGGING"`
	Mongo     config.MongoDB   `env:"MONGO"`

	Minio MinioConfig `env:"MINIO"`

	FilesBucket string `env:"FILES_BUCKET"`
}

type MinioConfig struct {
	Endpoint        string `env:"ENDPOINT"`
	AccessKeyID     string `env:"ACCESS_KEY_ID"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
	UseSSL          bool   `env:"USE_SSL"`
}

func defaultConfig() Config {
	return Config{
		Server: config.Server{
			Addr: ":4005",
		},

		Telemetry: telemetry.Config{
			Enable: true,
		},

		Logging: clog.Config{
			Format: clog.TextFormat,
			Debug:  true,
		},

		Mongo: config.MongoDB{
			URI:          "mongodb://localhost:27017",
			DatabaseName: "bloefish_svc_file_upload",
		},

		Minio: MinioConfig{
			Endpoint:        "localhost:9000",
			AccessKeyID:     "minio_key",
			SecretAccessKey: "minio_secret_key",
			UseSSL:          false,
		},

		FilesBucket: "bloefish-svc-files",
	}
}

func Run(ctx context.Context) error {
	cfg := defaultConfig()
	config.MustHydrateFromEnvironment(ctx, &cfg)

	shutdown := cfg.Telemetry.MustSetup(ctx)
	defer func() {
		if err := shutdown(ctx); err != nil {
			clog.Get(ctx).WithError(err).Error("failed to shutdown telemetry")
		}
	}()

	ctx = clog.Set(ctx, cfg.Logging.Configure(ctx))
	_, mongoDatabase := cfg.Mongo.MustConnect(ctx)

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
	if err != nil {
		return err
	}
	if err := ensureBucketExists(ctx, minioClient, cfg.FilesBucket); err != nil {
		return fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	app := &app.App{
		FileRepository:    repositories.NewMgoFile(mongoDatabase),
		FileObjectService: services.NewMinioFileObject(minioClient, cfg.FilesBucket),
	}

	rpc := rpc.New(ctx, app)

	return rpc.Run(ctx, cfg.Server)
}

func ensureBucketExists(ctx context.Context, client *minio.Client, bucketName string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	return nil
}
