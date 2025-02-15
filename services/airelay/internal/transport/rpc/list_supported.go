package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/airelay"
)

func (r *RPC) ListSupported(ctx context.Context) (*airelay.ListSupportedResponse, error) {
	return r.app.ListSupported(ctx)
}
