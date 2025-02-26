package app

import (
	"context"
	"fmt"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/forkedcontext"
	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/0xdeafcafe/bloefish/services/skillset"
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

	skillSets, err := a.SkillSetService.GetManySkillSets(ctx, &skillset.GetManySkillSetsRequest{
		SkillSetIDs: req.SkillSetIDs,
		Owner: &skillset.Actor{
			Type:       skillset.ActorType(req.Owner.Type),
			Identifier: req.Owner.Identifier,
		},
		AllowDeleted: false,
	})
	if err != nil {
		return nil, err
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
		SkillSetIDs:    req.SkillSetIDs,
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
		SkillSetIDs:    []string{},
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
			SkillSets:          skillSets.SkillSets,
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

			if saveErrorErr := a.InteractionRepository.AddError(ctx, activeInteraction.ID, cher.Coerce(err)); saveErrorErr != nil {
				clog.Get(ctx).WithError(saveErrorErr).Error("failed to save error to interaction")
			}

			clog.Get(ctx).WithError(err).Error("failed to create conversation message reply")
		}

		return nil
	})

	return &conversation.CreateConversationMessageResponse{
		ConversationID: convo.ID,
		InputInteraction: &conversation.CreateConversationMessageResponseInteraction{
			ID:              interaction.ID,
			FileIDs:         interaction.FileIDs,
			SkillSetIDs:     interaction.SkillSetIDs,
			StreamChannelID: streamingChannelID,

			MarkedAsExcludedAt: interaction.MarkedAsExcludedAt,

			MessageContent: interaction.MessageContent,
			Errors:         interaction.Errors,

			AIRelayOptions: &conversation.AIRelayOptions{
				ProviderID: interaction.AIRelayOptions.ProviderID,
				ModelID:    interaction.AIRelayOptions.ModelID,
			},
			Owner: &conversation.Actor{
				Type:       conversation.ActorType(interaction.Owner.Type),
				Identifier: interaction.Owner.Identifier,
			},

			CreatedAt:   interaction.CreatedAt,
			UpdatedAt:   interaction.UpdatedAt,
			CompletedAt: interaction.CompletedAt,
			DeletedAt:   interaction.DeletedAt,
		},
		ResponseInteraction: &conversation.CreateConversationMessageResponseInteraction{
			ID:              activeInteraction.ID,
			FileIDs:         activeInteraction.FileIDs,
			SkillSetIDs:     activeInteraction.SkillSetIDs,
			StreamChannelID: streamingChannelID,

			MarkedAsExcludedAt: activeInteraction.MarkedAsExcludedAt,

			MessageContent: activeInteraction.MessageContent,
			Errors:         activeInteraction.Errors,

			AIRelayOptions: &conversation.AIRelayOptions{
				ProviderID: activeInteraction.AIRelayOptions.ProviderID,
				ModelID:    activeInteraction.AIRelayOptions.ModelID,
			},
			Owner: &conversation.Actor{
				Type:       conversation.ActorType(activeInteraction.Owner.Type),
				Identifier: activeInteraction.Owner.Identifier,
			},

			CreatedAt:   activeInteraction.CreatedAt,
			UpdatedAt:   activeInteraction.UpdatedAt,
			CompletedAt: activeInteraction.CompletedAt,
			DeletedAt:   activeInteraction.DeletedAt,
		},
		StreamChannelID: streamingChannelID,
	}, nil
}
