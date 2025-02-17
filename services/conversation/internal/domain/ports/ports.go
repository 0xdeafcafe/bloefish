package ports

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation/internal/domain/models"
)

type ConversationRepository interface {
	Create(ctx context.Context, cmd *models.CreateConversationCommand) (*models.Conversation, error)
	GetByID(ctx context.Context, conversationID string) (*models.Conversation, error)
	ListByOwner(ctx context.Context, actor models.Actor) ([]*models.Conversation, error)
}

type InteractionRepository interface {
	Create(ctx context.Context, cmd *models.CreateInteractionCommand) (*models.Interaction, error)
	CreateActive(ctx context.Context, cmd *models.CreateActiveInteractionCommand) (*models.Interaction, error)
	MarkActiveAsComplete(ctx context.Context, interactionID, messageContent string) error
	GetByID(ctx context.Context, interactionID string) (*models.Interaction, error)
	GetAllByConversationID(ctx context.Context, conversationID string) ([]*models.Interaction, error)
}
