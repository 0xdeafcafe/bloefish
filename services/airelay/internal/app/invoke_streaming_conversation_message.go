package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/openai/openai-go"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/0xdeafcafe/bloefish/services/airelay/internal/libraries/relay"
	"github.com/0xdeafcafe/bloefish/services/stream"
)

func (a *App) InvokeStreamingConversationMessage(ctx context.Context, req *airelay.InvokeStreamingConversationMessageRequest) (*airelay.InvokeStreamingConversationMessageResponse, error) {
	fileIDs := []string{}
	for _, msg := range req.Messages {
		fileIDs = append(fileIDs, msg.FileIDs...)
	}

	downloadedFiles, err := a.downloadFiles(ctx, req.Owner, fileIDs)
	if err != nil {
		return nil, err
	}

	messages := make([]relay.Message, len(req.Messages))
	for i, msg := range req.Messages {
		fileContent := ""

		if len(msg.FileIDs) > 0 {
			for _, fileID := range msg.FileIDs {
				file := downloadedFiles[fileID]
				if file == nil {
					return nil, cher.New("file_missing_from_downloads", cher.M{
						"file_id":         fileID,
						"downloads_count": len(downloadedFiles),
						"message_index":   i,
					})
				}

				fileContent += fmt.Sprintf("\n\nFile name: %s\nFile content:\n%s", file.Name, string(file.Content))
			}
		}

		switch msg.Owner.Type {
		case airelay.ActorTypeBot:
			messages[i] = relay.Message{
				Role:    relay.RoleAssistant,
				Content: msg.Content + fileContent,
			}
		case airelay.ActorTypeUser:
			messages[i] = relay.Message{
				Role:    relay.RoleUser,
				Content: msg.Content + fileContent,
			}
		}
	}

	messages = append([]relay.Message{{
		Role:    relay.RoleSystem,
		Content: systemInstructionMessage,
	}}, messages...)

	chatStream, err := a.Relay.With(req.AIRelayOptions.ProviderID).NewChatStream(ctx, relay.ChatStreamParams{
		ModelID:  req.AIRelayOptions.ModelID,
		Messages: messages,
	})
	if err != nil {
		if errors.Is(err, relay.ErrRequiredProviderMissing) {
			return nil, cher.New("unsupported_ai_provider", cher.M{
				"provider_id": req.AIRelayOptions.ProviderID,
			})
		}

		return nil, err
	}

	for chatStream.Next() {
		event := chatStream.Current()

		if event.Content != "" {
			if err := a.StreamService.SendMessageFragment(ctx, &stream.SendMessageFragmentRequest{
				ChannelID:      req.StreamingChannelID,
				MessageContent: event.Content,
			}); err != nil {
				return nil, fmt.Errorf("failed to send message fragment: %w", err)
			}
		}
	}

	if err := chatStream.Err(); err != nil {
		var coercedError cher.E

		var apierr *openai.Error
		if errors.As(err, &apierr) {
			switch apierr.StatusCode {
			case http.StatusNotFound:
				coercedError = cher.New("ai_model_not_found", cher.M{
					"model_id": req.AIRelayOptions.ModelID,
				})
			}
		}
		if coercedError.Code == "" {
			coercedError = cher.Coerce(err)
		}

		if err := a.StreamService.SendErrorMessage(ctx, &stream.SendErrorMessageRequest{
			ChannelID: req.StreamingChannelID,
			Error:     coercedError,
		}); err != nil {
			return nil, fmt.Errorf("failed to send error message: %w", err)
		}

		return nil, coercedError
	}

	return &airelay.InvokeStreamingConversationMessageResponse{
		MessageContent: chatStream.Content(),
	}, nil
}
