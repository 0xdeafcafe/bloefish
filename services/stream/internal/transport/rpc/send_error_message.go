package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/stream"
)

func (r *RPC) SendErrorMessage(ctx context.Context, req *stream.SendErrorMessageRequest) error {
	return r.app.SendErrorMessage(ctx, req)
}
