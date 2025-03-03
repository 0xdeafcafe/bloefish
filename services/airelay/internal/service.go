package internal

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/libraries/telemetry"
	"github.com/0xdeafcafe/bloefish/services/airelay/internal/app"
	"github.com/0xdeafcafe/bloefish/services/airelay/internal/libraries/relay"
	"github.com/0xdeafcafe/bloefish/services/airelay/internal/libraries/relay/providers/ollama"
	"github.com/0xdeafcafe/bloefish/services/airelay/internal/libraries/relay/providers/openai"
	"github.com/0xdeafcafe/bloefish/services/airelay/internal/transport/rpc"
	"github.com/0xdeafcafe/bloefish/services/conversation"
	"github.com/0xdeafcafe/bloefish/services/fileupload"
	"github.com/0xdeafcafe/bloefish/services/stream"

	oaiClient "github.com/openai/openai-go"
	openaiOption "github.com/openai/openai-go/option"
)

type Config struct {
	Server    config.Server    `env:"SERVER"`
	Telemetry telemetry.Config `env:"TELEMETRY"`
	Logging   clog.Config      `env:"LOGGING"`

	ConversationService config.UnauthenticatedService `env:"CONVERSATION_SERVICE"`
	FileUploadService   config.UnauthenticatedService `env:"FILE_UPLOAD_SERVICE"`
	StreamService       config.UnauthenticatedService `env:"STREAM_SERVICE"`

	AIProviders AIProviders `env:"AI_PROVIDERS"`
}

type AIProviders struct {
	OpenAI OpenAIConfig `env:"OPENAI"`
	Ollama OllamaConfig `env:"OLLAMA"`
}

type OpenAIConfig struct {
	APIKey         string `env:"API_KEY"`
	OrganizationID string `env:"ORGANIZATION_ID"`
}

type OllamaConfig struct {
	Endpoint string `env:"ENDPOINT"`
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

		Telemetry: telemetry.Config{
			Enable: true,
		},

		AIProviders: AIProviders{
			OpenAI: OpenAIConfig{},
			Ollama: OllamaConfig{
				Endpoint: "http://localhost:11434",
			},
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

	app := &app.App{
		Relay: relay.NewClient(
			relay.WithProvider(openai.NewProvider(
				oaiClient.NewClient(
					openaiOption.WithAPIKey(cfg.AIProviders.OpenAI.APIKey),
				),
				openai.WithModels([]openai.Model{{
					ID:   string(oaiClient.ChatModelGPT4),
					Name: "GPT 4",
				}, {
					ID:   string(oaiClient.ChatModelGPT4Turbo),
					Name: "GPT 4 turbo",
				}, {
					ID:   string(oaiClient.ChatModelGPT4o),
					Name: "GPT 4o",
				}, {
					ID:   string(oaiClient.ChatModelGPT4oMini),
					Name: "GPT 4o mini",
				}, {
					ID:   string(oaiClient.ChatModelGPT3_5Turbo),
					Name: "GPT 3.5 turbo",
				}, {
					ID:   string(oaiClient.ChatModelO1Mini),
					Name: "o1 mini",
				}}),
			)),
			relay.WithProvider(ollama.NewProvider()),
		),

		ConversationService: conversation.NewRPCClient(ctx, cfg.ConversationService),
		FileUploadService:   fileupload.NewRPCClient(ctx, cfg.FileUploadService),
		StreamService:       stream.NewRPCClient(ctx, cfg.StreamService),
	}

	rpc := rpc.New(ctx, app)

	return rpc.Run(ctx, cfg.Server)
}
