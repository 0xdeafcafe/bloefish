package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/forkedcontext"
	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/0xdeafcafe/bloefish/services/stream"

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

	if convo.Title == nil {
		if newConversation, err := a.InteractionRepository.ConversationHasInteractions(ctx, convo.ID); err != nil {
			return nil, err
		} else if newConversation {
			titleStreamingChannelID := fmt.Sprintf("%s/title", convo.ID)

			forkedcontext.ForkContext(ctx, func(ctx context.Context) error {
				if err := a.generateConversationTitle(ctx, &generateConversationTitleCommand{
					Conversation: convo,
					Owner: &airelay.Actor{
						Type:       airelay.ActorType(req.Owner.Type),
						Identifier: req.Owner.Identifier,
					},
					Interaction:        interaction,
					StreamingChannelID: titleStreamingChannelID,
					UseStreaming:       req.Options.UseStreaming,
				}); err != nil {
					if sendErrorErr := a.StreamService.SendErrorMessage(ctx, &stream.SendErrorMessageRequest{
						ChannelID: titleStreamingChannelID,
						Error:     cher.Coerce(err),
					}); sendErrorErr != nil {
						clog.Get(ctx).WithError(sendErrorErr).Error("failed to send title error message to stream service")
					}

					clog.Get(ctx).WithError(err).Error("failed to generate conversation title")
				}

				return nil
			})
		}
	}

	activeInteraction, err := a.InteractionRepository.CreateActive(ctx, &models.CreateActiveInteractionCommand{
		IdempotencyKey: fmt.Sprintf("%s-response", req.IdempotencyKey),
		ConversationID: convo.ID,
		FileIDs:        []string{},
		MessageContent: "",
		Owner: &models.CreateActiveInteractionCommandOwner{
			Type:       models.ActorTypeBot,
			Identifier: interactionAIRelayOptions.ProviderID,
		},
		AIRelayOptions: &models.CreateActiveInteractionCommandAIRelayOptions{
			ProviderID: interactionAIRelayOptions.ProviderID,
			ModelID:    interactionAIRelayOptions.ModelID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create active interaction: %w", err)
	}

	streamingChannelID := fmt.Sprintf("%s/%s", convo.ID, activeInteraction.ID)
	forkedcontext.ForkContext(ctx, func(ctx context.Context) error {
		if err := a.createConversationMessageReply(ctx, &createConversationMessageReplyCommand{
			Conversation: convo,
			Owner: &airelay.Actor{
				Type:       airelay.ActorType(interaction.Owner.Type),
				Identifier: interaction.Owner.Identifier,
			},
			Interaction:        interaction,
			ActiveInteraction:  activeInteraction,
			StreamingChannelID: streamingChannelID,
			UseStreaming:       req.Options.UseStreaming,
		}); err != nil {
			if sendErrorErr := a.StreamService.SendErrorMessage(ctx, &stream.SendErrorMessageRequest{
				ChannelID: streamingChannelID,
				Error:     cher.Coerce(err),
			}); sendErrorErr != nil {
				clog.Get(ctx).WithError(sendErrorErr).Error("failed to send error message to stream service")
			}

			clog.Get(ctx).WithError(err).Error("failed to create conversation message reply")
		}

		return nil
	})

	return &conversation.CreateConversationMessageResponse{
		ConversationID: convo.ID,
		InputInteraction: &conversation.CreateConversationMessageResponseInteraction{
			ID: interaction.ID,
			Owner: &conversation.Actor{
				Type:       conversation.ActorType(interaction.Owner.Type),
				Identifier: interaction.Owner.Identifier,
			},
			FileIDs:        interaction.FileIDs,
			MessageContent: interaction.MessageContent,
			AIRelayOptions: &conversation.AIRelayOptions{
				ProviderID: interaction.AIRelayOptions.ProviderID,
				ModelID:    interaction.AIRelayOptions.ModelID,
			},
			CreatedAt:   interaction.CreatedAt,
			UpdatedAt:   interaction.UpdatedAt,
			CompletedAt: interaction.CompletedAt,
			DeletedAt:   interaction.DeletedAt,
		},
		ResponseInteraction: &conversation.CreateConversationMessageResponseInteraction{
			ID: activeInteraction.ID,
			Owner: &conversation.Actor{
				Type:       conversation.ActorType(activeInteraction.Owner.Type),
				Identifier: activeInteraction.Owner.Identifier,
			},
			FileIDs:        activeInteraction.FileIDs,
			MessageContent: activeInteraction.MessageContent,
			AIRelayOptions: &conversation.AIRelayOptions{
				ProviderID: activeInteraction.AIRelayOptions.ProviderID,
				ModelID:    activeInteraction.AIRelayOptions.ModelID,
			},
			StreamChannelID: streamingChannelID,
			CreatedAt:       activeInteraction.CreatedAt,
			UpdatedAt:       activeInteraction.UpdatedAt,
			CompletedAt:     activeInteraction.CompletedAt,
			DeletedAt:       activeInteraction.DeletedAt,
		},
		StreamChannelID: streamingChannelID,
	}, nil
}

type createConversationMessageReplyCommand struct {
	Conversation *models.Conversation
	Owner        *airelay.Actor

	Interaction       *models.Interaction
	ActiveInteraction *models.Interaction

	StreamingChannelID string
	UseStreaming       bool
}

func (a *App) createConversationMessageReply(
	ctx context.Context,
	cmd *createConversationMessageReplyCommand,
) error {
	conversationInteractions, err := a.InteractionRepository.GetAllByConversationID(ctx, cmd.Conversation.ID)
	if err != nil {
		return err
	}

	messages := make([]*airelay.InvokeConversationMessageRequestMessage, 0, len(conversationInteractions))

	for _, interaction := range conversationInteractions {
		if interaction.CompletedAt == nil || interaction.MarkedAsExcludedAt != nil {
			continue
		}

		messages = append(messages, &airelay.InvokeConversationMessageRequestMessage{
			Owner: &airelay.Actor{
				Type:       airelay.ActorType(interaction.Owner.Type),
				Identifier: interaction.Owner.Identifier,
			},
			FileIDs: []string{}, // TODO(afr): Handle files
			Content: interaction.MessageContent,
		})
	}

	aiRelayOptions := &models.AIRelayOptions{
		ProviderID: cmd.Conversation.AIRelayOptions.ProviderID,
		ModelID:    cmd.Conversation.AIRelayOptions.ModelID,
	}
	if cmd.Interaction.AIRelayOptions != nil {
		aiRelayOptions = &models.AIRelayOptions{
			ProviderID: cmd.Interaction.AIRelayOptions.ProviderID,
			ModelID:    cmd.Interaction.AIRelayOptions.ModelID,
		}
	}

	var messageContent string
	if cmd.UseStreaming {
		response, err := a.AIRelayService.InvokeStreamingConversationMessage(ctx, &airelay.InvokeStreamingConversationMessageRequest{
			StreamingChannelID: cmd.StreamingChannelID,
			Owner:              cmd.Owner,
			Messages:           messages,
			AIRelayOptions: &airelay.InvokeConversationMessageRequestAIRelayOptions{
				ProviderID: aiRelayOptions.ProviderID,
				ModelID:    aiRelayOptions.ModelID,
			},
		})
		if err != nil {
			json, _ := json.Marshal(err)
			clog.Get(ctx).WithError(err).WithField("json", string(json)).Error("failed to invoke streaming conversation message")

			return err
		}

		messageContent = response.MessageContent
	} else {
		response, err := a.AIRelayService.InvokeConversationMessage(ctx, &airelay.InvokeConversationMessageRequest{
			Owner:    cmd.Owner,
			Messages: messages,
			AIRelayOptions: &airelay.InvokeConversationMessageRequestAIRelayOptions{
				ProviderID: aiRelayOptions.ProviderID,
				ModelID:    aiRelayOptions.ModelID,
			},
		})
		if err != nil {
			json, _ := json.Marshal(err)
			clog.Get(ctx).WithError(err).WithField("json", string(json)).Error("failed to invoke conversation message")

			return err
		}

		messageContent = response.MessageContent
	}

	if err := a.InteractionRepository.MarkActiveAsComplete(ctx, cmd.ActiveInteraction.ID, messageContent); err != nil {
		return err
	}

	return nil
}
