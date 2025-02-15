package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/airelay"
)

func (r *RPC) InvokeStreamingConversationMessage(ctx context.Context, req *airelay.InvokeStreamingConversationMessageRequest) (*airelay.InvokeStreamingConversationMessageResponse, error) {
	return r.app.InvokeStreamingConversationMessage(ctx, req)
}
