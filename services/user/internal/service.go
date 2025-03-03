package internal

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/libraries/telemetry"
	"github.com/0xdeafcafe/bloefish/services/user/internal/app"
	"github.com/0xdeafcafe/bloefish/services/user/internal/app/repositories"
	"github.com/0xdeafcafe/bloefish/services/user/internal/transport/rpc"
)

type Config struct {
	Server    config.Server    `env:"SERVER"`
	Telemetry telemetry.Config `env:"TELEMETRY"`
	Logging   clog.Config      `env:"LOGGING"`
	Mongo     config.MongoDB   `env:"MONGO"`
}

func defaultConfig() Config {
	return Config{
		Server: config.Server{
			Addr: ":4001",
		},

		Logging: clog.Config{
			Format: clog.TextFormat,
			Debug:  true,
		},
		Telemetry: telemetry.Config{
			Enable: true,
		},

		Mongo: config.MongoDB{
			URI:          "mongodb://localhost:27017",
			DatabaseName: "bloefish_svc_user",
		},
	}
}

func Run(ctx context.Context) error {
	cfg := defaultConfig()
	config.MustHydrate(ctx, &cfg)

	shutdown := cfg.Telemetry.MustSetup(ctx)
	defer func() {
		if err := shutdown(ctx); err != nil {
			clog.Get(ctx).WithError(err).Error("failed to shutdown telemetry")
		}
	}()

	ctx = clog.Set(ctx, cfg.Logging.Configure(ctx))
	_, mongoDatabase := cfg.Mongo.MustConnect(ctx)

	app := &app.App{
		UserRepository: repositories.NewMgoUser(mongoDatabase),
	}

	rpc := rpc.New(ctx, app)

	return rpc.Run(ctx, cfg.Server)
}
