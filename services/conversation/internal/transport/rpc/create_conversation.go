package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation"
)

func (r *RPC) CreateConversation(ctx context.Context, req *conversation.CreateConversationRequest) (*conversation.CreateConversationResponse, error) {
	return r.app.CreateConversation(ctx, req)
}
