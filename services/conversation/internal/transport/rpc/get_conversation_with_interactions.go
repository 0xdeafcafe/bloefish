package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation"
)

func (r *RPC) GetConversationWithInteractions(ctx context.Context, req *conversation.GetConversationWithInteractionsRequest) (*conversation.GetConversationWithInteractionsResponse, error) {
	return r.app.GetConversationWithInteractions(ctx, req)
}
