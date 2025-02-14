package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/stream"
)

func (r *RPC) HandleSendMessageFragment(ctx context.Context, req *stream.SendMessageFragmentRequest) error {
	return r.app.SendMessageFragment(ctx, req)
}
