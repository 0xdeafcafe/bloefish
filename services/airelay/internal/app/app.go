package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/0xdeafcafe/bloefish/services/conversation"
	"github.com/0xdeafcafe/bloefish/services/fileupload"
	"github.com/0xdeafcafe/bloefish/services/stream"
	"github.com/openai/openai-go"
)

type App struct {
	OpenAI          *openai.Client
	SupportedModels []openai.ChatModel

	ConversationService conversation.Service
	FileUploadService   fileupload.Service
	StreamService       stream.Service
}

func (a *App) ListSupported(ctx context.Context) (*airelay.ListSupportedResponse, error) {
	if a.OpenAI == nil {
		return &airelay.ListSupportedResponse{}, nil
	}

	return &airelay.ListSupportedResponse{
		Providers: []*airelay.ListSupportedResponseProvider{{
			ID:   "provider_open_ai",
			Name: "Open AI",
			Models: []*airelay.ListSupportedResponseProviderModel{{
				ID: openai.ChatModelGPT4,
			}, {
				ID: openai.ChatModelGPT4o,
			}, {
				ID: openai.ChatModelO1,
			}, {
				ID: openai.ChatModelO1Mini,
			}, {
				ID: openai.ChatModelO3Mini,
			}},
		}},
	}, nil
}
