package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/0xdeafcafe/bloefish/services/skillset"

	"github.com/0xdeafcafe/bloefish/services/conversation/internal/domain/models"
)

type createConversationMessageReplyCommand struct {
	Conversation *models.Conversation
	Owner        *airelay.Actor

	Interaction       *models.Interaction
	ActiveInteraction *models.Interaction

	SkillSets []*skillset.SkillSet

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

	// Handle skill sets
	if len(cmd.SkillSets) > 0 {
		for _, skillSet := range cmd.SkillSets {
			messages = append(messages, &airelay.InvokeConversationMessageRequestMessage{
				Owner: &airelay.Actor{
					Type:       airelay.ActorType(skillSet.Owner.Type),
					Identifier: skillSet.Owner.Identifier,
				},
				FileIDs: []string{}, // Skill sets current can't have files
				Content: fmt.Sprintf(
					"Use the following instructions to guide your responses or to learn more context about the subject: %s",
					skillSet.Prompt,
				),
			})
		}
	}

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
