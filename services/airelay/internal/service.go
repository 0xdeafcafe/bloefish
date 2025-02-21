package internal

import (
	"context"

	"github.com/openai/openai-go"
	oai_option "github.com/openai/openai-go/option"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/libraries/ollama"
	"github.com/0xdeafcafe/bloefish/services/airelay/internal/app"
	"github.com/0xdeafcafe/bloefish/services/airelay/internal/libraries/relay"
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
	config.MustHydrateFromEnvironment(ctx, &cfg)
	ctx = clog.Set(ctx, cfg.Logging.Configure(ctx))

	app := &app.App{
		Relay: relay.NewClient(
			relay.WithOllamaClient(ollama.NewClient(
				ollama.WithEndpointURL(cfg.AIProviders.Ollama.Endpoint),
			)),
			relay.WithOpenAIClient(openai.NewClient(
				oai_option.WithAPIKey(cfg.AIProviders.OpenAI.APIKey),
			)),
			relay.WithProviderModels([]relay.Provider{{
				ID:   "open_ai",
				Name: "Open AI",
				Models: []relay.Model{{
					ID:   openai.ChatModelGPT4,
					Name: "GPT 4",
				}, {
					ID:   openai.ChatModelGPT4Turbo,
					Name: "GPT 4 turbo",
				}, {
					ID:   openai.ChatModelGPT4o,
					Name: "GPT 4o",
				}, {
					ID:   openai.ChatModelGPT4oMini,
					Name: "GPT 4o mini",
				}, {
					ID:   openai.ChatModelGPT3_5Turbo,
					Name: "GPT 3.5 turbo",
				}, {
					ID:   openai.ChatModelO1Mini,
					Name: "o1 mini",
				}},
			}, {
				ID:   "ollama",
				Name: "Ollama",
				Models: []relay.Model{{
					ID:   "phi3:mini",
					Name: "Phi3 mini",
				}},
			}}),
		),

		ConversationService: conversation.NewRPCClient(ctx, cfg.ConversationService),
		FileUploadService:   fileupload.NewRPCClient(ctx, cfg.FileUploadService),
		StreamService:       stream.NewRPCClient(ctx, cfg.StreamService),
	}

	rpc := rpc.New(ctx, app)

	return rpc.Run(ctx, cfg.Server)
}
