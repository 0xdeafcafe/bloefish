package internal

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/services/airelay/internal/app"
	"github.com/0xdeafcafe/bloefish/services/airelay/internal/transport/rpc"
	"github.com/0xdeafcafe/bloefish/services/conversation"
	"github.com/0xdeafcafe/bloefish/services/fileupload"
	"github.com/0xdeafcafe/bloefish/services/stream"
)

type Config struct {
	Server  config.Server `env:"SERVER"`
	Logging clog.Config   `env:"LOGGING"`

	ConversationService config.UnauthenticatedService `env:"CONVERSATION_SERVICE"`
	FileUploadService   config.UnauthenticatedService `env:"FILE_UPLOAD_SERVICE"`
	StreamService       config.UnauthenticatedService `env:"STREAM_SERVICE"`

	AIProviders AIProviders `env:"AI_PROVIDERS"`
}

type AIProviders struct {
	OpenAI OpenAIConfig `env:"OPENAI"`
}

type OpenAIConfig struct {
	APIKey         string `env:"API_KEY"`
	OrganizationID string `env:"ORGANIZATION_ID"`
}

func defaultConfig() Config {
	return Config{
		Server: config.Server{
			Addr: ":4003",
		},

		Logging: clog.Config{
			Format: clog.TextFormat,
			Debug:  true,
		},

		ConversationService: config.UnauthenticatedService{
			BaseURL: "http://localhost:4002/rpc",
		},
		FileUploadService: config.UnauthenticatedService{
			BaseURL: "http://localhost:4005/rpc",
		},
		StreamService: config.UnauthenticatedService{
			BaseURL: "http://localhost:4004/rpc",
		},

		AIProviders: AIProviders{
			OpenAI: OpenAIConfig{},
		},
	}
}

func Run(ctx context.Context) error {
	cfg := defaultConfig()
	config.MustHydrateFromEnvironment(ctx, &cfg)
	ctx = clog.Set(ctx, cfg.Logging.Configure(ctx))

	app := &app.App{
		ConversationService: conversation.NewRPCClient(ctx, cfg.ConversationService),
		FileUploadService:   fileupload.NewRPCClient(ctx, cfg.FileUploadService),
		StreamService:       stream.NewRPCClient(ctx, cfg.StreamService),

		OpenAI: openai.NewClient(
			option.WithAPIKey(cfg.AIProviders.OpenAI.APIKey),
		),
		SupportedModels: []openai.ChatModel{},
	}

	rpc := rpc.New(ctx, app)

	return rpc.Run(ctx, cfg.Server)
}
