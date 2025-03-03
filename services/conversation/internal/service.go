package internal

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/libraries/telemetry"
	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/app"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/app/repositories"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/transport/rpc"
	"github.com/0xdeafcafe/bloefish/services/skillset"
	"github.com/0xdeafcafe/bloefish/services/stream"
	"github.com/0xdeafcafe/bloefish/services/user"
)

type Config struct {
	Server    config.Server    `env:"SERVER"`
	Telemetry telemetry.Config `env:"TELEMETRY"`
	Logging   clog.Config      `env:"LOGGING"`
	Mongo     config.MongoDB   `env:"MONGO"`

	AIRelayService  config.UnauthenticatedService `env:"AI_RELAY_SERVICE"`
	SkillSetService config.UnauthenticatedService `env:"SKILL_SET_SERVICE"`
	StreamService   config.UnauthenticatedService `env:"STREAM_SERVICE"`
	UserService     config.UnauthenticatedService `env:"USER_SERVICE"`
}

func defaultConfig() Config {
	return Config{
		Server: config.Server{
			Addr: ":4002",
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
			DatabaseName: "bloefish_svc_conversation",
		},

		AIRelayService: config.UnauthenticatedService{
			BaseURL: "http://localhost:4003/rpc",
		},
		SkillSetService: config.UnauthenticatedService{
			BaseURL: "http://localhost:4006/rpc",
		},
		StreamService: config.UnauthenticatedService{
			BaseURL: "http://localhost:4004/rpc",
		},
		UserService: config.UnauthenticatedService{
			BaseURL: "http://localhost:4001/rpc",
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
		ConversationRepository: repositories.NewMgoConversation(mongoDatabase),
		InteractionRepository:  repositories.NewMgoInteraction(mongoDatabase),

		AIRelayService:  airelay.NewRPCClient(ctx, cfg.AIRelayService),
		SkillSetService: skillset.NewRPCClient(ctx, cfg.SkillSetService),
		StreamService:   stream.NewRPCClient(ctx, cfg.StreamService),
		UserService:     user.NewRPCClient(ctx, cfg.UserService),
	}

	rpc := rpc.New(ctx, app)

	return rpc.Run(ctx, cfg.Server)
}
