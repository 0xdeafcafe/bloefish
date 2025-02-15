package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/user"
)

func (r *RPC) GetOrCreateDefaultUser(ctx context.Context) (*user.GetOrCreateDefaultUserResponse, error) {
	userResponse, err := r.app.GetOrCreateDefaultUser(ctx)
	if err != nil {
		return nil, err
	}

	return &user.GetOrCreateDefaultUserResponse{
		User: userResponse,
	}, nil
}
