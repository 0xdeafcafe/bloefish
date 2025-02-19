package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/domain/models"
)

type generateConversationTitleCommand struct {
	Conversation *models.Conversation
	Owner        *airelay.Actor

	Interaction *models.Interaction

	StreamingChannelID string
	UseStreaming       bool
}

func (a *App) generateConversationTitle(
	ctx context.Context,
	cmd *generateConversationTitleCommand,
) error {
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
	messages := []*airelay.InvokeConversationMessageRequestMessage{{
		Content: fmt.Sprintf(
			"Generate a conversation title for the following text. It should be concise but descriptive. It should be a single sentence, no more than 100 characters. Return nothing but the title. The output should not be wrapped in double quotes at all, or use any markdown formatting. Again, do not wrap the output in quotes, just return it as is.\n\n %s",
			cmd.Interaction.MessageContent,
		),
		Owner:   cmd.Owner,
		FileIDs: []string{},
	}}

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
			clog.Get(ctx).WithError(err).WithField("json", string(json)).Error("failed to call invoke streaming conversation message to generate title")

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
			clog.Get(ctx).WithError(err).WithField("json", string(json)).Error("failed to call invoke conversation message to generate title")

			return err
		}

		messageContent = response.MessageContent
	}

	if err := a.ConversationRepository.UpdateTitle(ctx, cmd.Conversation.ID, messageContent); err != nil {
		return fmt.Errorf("failed to update conversation title: %w", err)
	}

	return nil
}
