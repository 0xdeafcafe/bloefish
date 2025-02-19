package app

import (
	"context"
	"runtime"

	"github.com/0xdeafcafe/bloefish/services/conversation"
	"golang.org/x/sync/errgroup"
)

func (a *App) DeleteConversations(ctx context.Context, req *conversation.DeleteConversationsRequest) error {
	errGroup, egCtx := errgroup.WithContext(ctx)
	errGroup.SetLimit(runtime.NumCPU() * 4)

	errGroup.Go(func() error {
		return a.ConversationRepository.DeleteMany(egCtx, req.ConversationIDs)
	})

	for _, conversationID := range req.ConversationIDs {
		errGroup.Go(func() error {
			return a.InteractionRepository.DeleteManyByConversationID(egCtx, conversationID)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return err
	}

	return nil
}
