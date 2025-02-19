package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation"
)

func (a *App) DeleteInteractions(ctx context.Context, req *conversation.DeleteInteractionsRequest) error {
	if err := a.ConversationRepository.DeleteMany(ctx, req.InteractionIDs); err != nil {
		return err
	}

	return nil
}
