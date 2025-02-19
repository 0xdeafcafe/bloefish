package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/conversation"
)

func (r *RPC) UpdateInteractionExcludedState(ctx context.Context, req *conversation.UpdateInteractionExcludedStateRequest) error {
	return r.app.UpdateInteractionExcludedState(ctx, req)
}
