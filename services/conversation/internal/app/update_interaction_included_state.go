package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation"
)

func (a *App) UpdateInteractionExcludedState(ctx context.Context, req *conversation.UpdateInteractionExcludedStateRequest) error {
	if err := a.InteractionRepository.UpdateExcludedState(ctx, req.InteractionID, req.Excluded); err != nil {
		return err
	}

	return nil
}
