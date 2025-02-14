package app

import (
	"context"
	"fmt"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/forkedcontext"
	"github.com/0xdeafcafe/bloefish/services/airelay"

	"github.com/0xdeafcafe/bloefish/services/conversation"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/domain/models"
)

func (a *App) CreateConversationMessage(ctx context.Context, req *conversation.CreateConversationMessageRequest) (*conversation.CreateConversationMessageResponse, error) {
	defaultUser, err := a.UserService.GetOrCreateDefaultUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.Owner.Identifier != defaultUser.User.ID {
		return nil, cher.New("invalid_owner", cher.M{"identifier": req.Owner.Identifier})
	}

	convo, err := a.ConversationRepository.GetByID(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}
	if convo.Owner.Type != models.ActorType(req.Owner.Type) || convo.Owner.Identifier != req.Owner.Identifier {
		return nil, cher.New("invalid_owner", cher.M{
			"type":       req.Owner.Type,
			"identifier": req.Owner.Identifier,
		})
	}

	var interactionAIRelayOptions *models.CreateInteractionCommandAIRelayOptions
	if req.AIRelayOptions != nil {
		interactionAIRelayOptions = &models.CreateInteractionCommandAIRelayOptions{
			ProviderID: req.AIRelayOptions.ProviderID,
			ModelID:    req.AIRelayOptions.ModelID,
		}
	}

	interaction, err := a.InteractionRepository.Create(ctx, &models.CreateInteractionCommand{
		IdempotencyKey: req.IdempotencyKey,
		ConversationID: convo.ID,
		FileIDs:        req.FileIDs,
		MessageContent: req.MessageContent,
		Owner: &models.CreateInteractionCommandOwner{
			Type:       models.ActorType(req.Owner.Type),
			Identifier: req.Owner.Identifier,
		},
		AIRelayOptions: interactionAIRelayOptions,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create interaction: %w", err)
	}

	forkedcontext.ForkContext(ctx, func(ctx context.Context) error {
		conversationInteractions, err := a.InteractionRepository.GetAllByConversationID(ctx, convo.ID)
		if err != nil {
			return fmt.Errorf("failed to get full conversation: %w", err)
		}

		messages := make([]*airelay.InvokeConversationMessageRequestMessage, len(conversationInteractions))

		for i, interaction := range conversationInteractions {
			messages[i] = &airelay.InvokeConversationMessageRequestMessage{
				Owner: &airelay.Actor{
					Type:       airelay.ActorType(interaction.Owner.Type),
					Identifier: interaction.Owner.Identifier,
				},
				FileIDs: []string{},
				Content: interaction.MessageContent,
			}
		}

		var messageContent string

		aiRelayOptions := &airelay.InvokeConversationMessageRequestAIRelayOptions{
			ProviderID: convo.AIRelayOptions.ProviderID,
			ModelID:    convo.AIRelayOptions.ModelID,
		}
		if interaction.AIRelayOptions != nil {
			aiRelayOptions = &airelay.InvokeConversationMessageRequestAIRelayOptions{
				ProviderID: interaction.AIRelayOptions.ProviderID,
				ModelID:    interaction.AIRelayOptions.ModelID,
			}
		}

		if req.Options.UseStreaming {
			response, err := a.AIRelayService.InvokeStreamingConversationMessage(ctx, &airelay.InvokeStreamingConversationMessageRequest{
				StreamingChannelID: convo.ID,
				Owner: &airelay.Actor{
					Type:       airelay.ActorType(req.Owner.Type),
					Identifier: req.Owner.Identifier,
				},
				Messages:       messages,
				AIRelayOptions: aiRelayOptions,
			})
			if err != nil {
				clog.Get(ctx).WithError(err).Error("failed to invoke streaming conversation message")

				return fmt.Errorf("failed to invoke streaming conversation message: %w", err)
			}

			messageContent = response.MessageContent
		} else {
			response, err := a.AIRelayService.InvokeConversationMessage(ctx, &airelay.InvokeConversationMessageRequest{
				Owner: &airelay.Actor{
					Type:       airelay.ActorType(req.Owner.Type),
					Identifier: req.Owner.Identifier,
				},
				Messages: messages,
				AIRelayOptions: &airelay.InvokeConversationMessageRequestAIRelayOptions{
					ProviderID: req.AIRelayOptions.ProviderID,
					ModelID:    req.AIRelayOptions.ModelID,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to invoke conversation message: %w", err)
			}

			messageContent = response.MessageContent
		}

		if _, err = a.InteractionRepository.Create(ctx, &models.CreateInteractionCommand{
			IdempotencyKey: time.Now().Format(time.RFC3339Nano), // TODO(afr): Replace this
			ConversationID: convo.ID,
			FileIDs:        []string{},
			MessageContent: messageContent,
			Owner: &models.CreateInteractionCommandOwner{
				Type:       models.ActorTypeBot,
				Identifier: "openai",
			},
			AIRelayOptions: &models.CreateInteractionCommandAIRelayOptions{
				ProviderID: req.AIRelayOptions.ProviderID,
				ModelID:    req.AIRelayOptions.ModelID,
			},
		}); err != nil {
			return fmt.Errorf("failed to create response interaction: %w", err)
		}

		return nil
	})

	return &conversation.CreateConversationMessageResponse{
		ConversationID:  convo.ID,
		InteractionID:   interaction.ID,
		StreamChannelID: convo.ID,
	}, nil
}
