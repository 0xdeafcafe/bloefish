package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/airelay"
)

func (r *RPC) HandleListSupported(ctx context.Context) (*airelay.ListSupportedResponse, error) {
	return r.app.ListSupported(ctx)
}
