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
	messages := make([]relay.Message, len(req.Messages))
	for i, msg := range req.Messages {
		switch msg.Owner.Type {
		case airelay.ActorTypeBot:
			messages[i] = relay.Message{
				Role:    relay.RoleAssistant,
				Content: msg.Content,
			}
		case airelay.ActorTypeUser:
			messages[i] = relay.Message{
				Role:    relay.RoleUser,
				Content: msg.Content,
			}
		}
	}

	messages = append([]relay.Message{{
		Role:    relay.RoleSystem,
		Content: systemInstructionMessage,
	}}, messages...)

	chatStream, err := a.Relay.NewChatStream(ctx, relay.ChatStreamParams{
		ProviderID: req.AIRelayOptions.ProviderID,
		ModelID:    req.AIRelayOptions.ModelID,
		Messages:   messages,
	})
	if err != nil {
		if errors.Is(err, relay.ErrUnsupportedProvider) {
			return nil, cher.New("unsupported_provider", cher.M{
				"provider_id": req.AIRelayOptions.ProviderID,
			})
		}
		return nil, fmt.Errorf("failed to create chat stream: %w", err)
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
