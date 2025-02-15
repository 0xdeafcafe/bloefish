package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation"
)

func (r *RPC) CreateConversationMessage(ctx context.Context, req *conversation.CreateConversationMessageRequest) (*conversation.CreateConversationMessageResponse, error) {
	return r.app.CreateConversationMessage(ctx, req)
}
