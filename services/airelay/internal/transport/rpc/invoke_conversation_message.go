package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/airelay"
)

func (r *RPC) HandleInvokeConversationMessage(ctx context.Context, req *airelay.InvokeConversationMessageRequest) (*airelay.InvokeConversationMessageResponse, error) {
	return r.app.InvokeConversationMessage(ctx, req)
}
