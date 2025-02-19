package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation"
)

func (r *RPC) DeleteInteractions(ctx context.Context, req *conversation.DeleteInteractionsRequest) error {
	return r.app.DeleteInteractions(ctx, req)
}
