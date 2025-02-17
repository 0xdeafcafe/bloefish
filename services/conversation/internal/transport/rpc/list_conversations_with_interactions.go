package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation"
)

func (r *RPC) ListConversationsWithInteractions(ctx context.Context, req *conversation.ListConversationsWithInteractionsRequest) (*conversation.ListConversationsWithInteractionsResponse, error) {
	return r.app.ListConversationsWithInteractions(ctx, req)
}
