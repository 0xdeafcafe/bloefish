package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation"
)

func (r *RPC) DeleteConversations(ctx context.Context, req *conversation.DeleteConversationsRequest) error {
	return r.app.DeleteConversations(ctx, req)
}
