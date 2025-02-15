package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/stream"
)

func (r *RPC) SendMessageFull(ctx context.Context, req *stream.SendMessageFullRequest) error {
	return r.app.SendMessageFull(ctx, req)
}
