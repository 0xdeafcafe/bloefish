package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/services/conversation"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/domain/models"
)

func (a *App) CreateConversation(ctx context.Context, req *conversation.CreateConversationRequest) (*conversation.CreateConversationResponse, error) {
	defaultUser, err := a.UserService.GetOrCreateDefaultUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.Owner.Identifier != defaultUser.User.ID {
		return nil, cher.New("invalid_owner", cher.M{"identifier": req.Owner.Identifier})
	}

	convo, err := a.ConversationRepository.Create(ctx, &models.CreateConversationCommand{
		IdempotencyKey: req.IdempotencyKey,
		Owner: &models.CreateConversationCommandOwner{
			Type:       models.ActorType(req.Owner.Type),
			Identifier: req.Owner.Identifier,
		},
		AIRelayOptions: &models.CreateConversationCommandAIRelayOptions{
			ProviderID: req.AIRelayOptions.ProviderID,
			ModelID:    req.AIRelayOptions.ModelID,
		},
	})
	if err != nil {
		return nil, err
	}

	return &conversation.CreateConversationResponse{
		ID: convo.ID,
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
		CreatedAt:       convo.CreatedAt,
		UpdatedAt:       convo.UpdatedAt,
		DeletedAt:       convo.DeletedAt,
	}, nil
}
