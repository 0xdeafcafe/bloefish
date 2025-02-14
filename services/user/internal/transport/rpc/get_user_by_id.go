package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/user"
)

func (r *RPC) HandleGetUserByID(ctx context.Context, req *user.GetUserByIDRequest) (*user.GetUserByIDResponse, error) {
	userResponse, err := r.app.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &user.GetUserByIDResponse{
		User: userResponse,
	}, nil
}
