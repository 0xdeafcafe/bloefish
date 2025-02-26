package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation"
)

func (a *App) GetInteraction(ctx context.Context, req *conversation.GetInteractionRequest) (*conversation.GetInteractionResponse, error) {
	foundInteraction, err := a.InteractionRepository.GetByID(ctx, req.InteractionID)
	if err != nil {
		return nil, err
	}

	return &conversation.GetInteractionResponse{
		ID:             foundInteraction.ID,
		ConversationID: foundInteraction.ConversationID,
		FileIDs:        foundInteraction.FileIDs,
		SkillSetIDs:    foundInteraction.SkillSetIDs,

		MarkedAsExcludedAt: foundInteraction.MarkedAsExcludedAt,
		MessageContent:     foundInteraction.MessageContent,
		Errors:             foundInteraction.Errors,

		Owner: &conversation.Actor{
			Type:       conversation.ActorType(foundInteraction.Owner.Type),
			Identifier: foundInteraction.Owner.Identifier,
		},
		AIRelayOptions: &conversation.AIRelayOptions{
			ProviderID: foundInteraction.AIRelayOptions.ProviderID,
			ModelID:    foundInteraction.AIRelayOptions.ModelID,
		},

		CreatedAt:   foundInteraction.CreatedAt,
		UpdatedAt:   foundInteraction.UpdatedAt,
		DeletedAt:   foundInteraction.DeletedAt,
		CompletedAt: foundInteraction.CompletedAt,
	}, nil
}
