package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation"
)

func (r *RPC) GetInteraction(ctx context.Context, req *conversation.GetInteractionRequest) (*conversation.GetInteractionResponse, error) {
	return r.app.GetInteraction(ctx, req)
}
