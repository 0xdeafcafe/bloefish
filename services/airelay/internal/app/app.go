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
			ID:   "open_ai",
			Name: "Open AI",
			Models: []*airelay.ListSupportedResponseProviderModel{{
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
			// }, {
			// 	ID:   openai.ChatModelO1,
			// 	Name: "o1",
			// }, {
			// 	ID:   openai.ChatModelO1Mini,
			// 	Name: "o1 mini",
			// }, {
			// 	ID:   openai.ChatModelO3Mini,
			// 	Name: "o3 mini",
			// }},
		}},
	}, nil
}
