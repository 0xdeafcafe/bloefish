package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation"
)

func (a *App) GetConversationWithInteractions(ctx context.Context, req *conversation.GetConversationWithInteractionsRequest) (*conversation.GetConversationWithInteractionsResponse, error) {
	convo, err := a.ConversationRepository.GetByID(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}

	interactions, err := a.InteractionRepository.GetAllByConversationID(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}

	resp := &conversation.GetConversationWithInteractionsResponse{
		ConversationID: convo.ID,
		Owner: &conversation.Actor{
			Type:       conversation.ActorType(convo.Owner.Type),
			Identifier: convo.Owner.Identifier,
		},
		AIRelayOptions: &conversation.AIRelayOptions{
			ProviderID: convo.AIRelayOptions.ProviderID,
			ModelID:    convo.AIRelayOptions.ModelID,
		},
		Title:           convo.Title,
		StreamChannelID: convo.ID,

		Interactions: make([]*conversation.GetConversationWithInteractionsResponseInteraction, len(interactions)),

		CreatedAt: convo.CreatedAt,
		UpdatedAt: convo.UpdatedAt,
		DeletedAt: convo.DeletedAt,
	}

	for i, interaction := range interactions {
		resp.Interactions[i] = &conversation.GetConversationWithInteractionsResponseInteraction{
			ID: interaction.ID,
			Owner: &conversation.Actor{
				Type:       conversation.ActorType(interaction.Owner.Type),
				Identifier: interaction.Owner.Identifier,
			},
			MessageContent: interaction.MessageContent,
			FileIDs:        interaction.FileIDs,
			AIRelayOptions: &conversation.AIRelayOptions{
				ProviderID: interaction.AIRelayOptions.ProviderID,
				ModelID:    interaction.AIRelayOptions.ModelID,
			},
			CreatedAt:   interaction.CreatedAt,
			UpdatedAt:   interaction.UpdatedAt,
			CompletedAt: interaction.CompletedAt,
			DeletedAt:   interaction.DeletedAt,
		}
	}

	return resp, nil
}
