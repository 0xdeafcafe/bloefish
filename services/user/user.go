package user

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/user/internal/domain/models"
)

type Service interface {
	GetUserByID(ctx context.Context, req *GetUserByIDRequest) (*GetUserByIDResponse, error)
	GetOrCreateDefaultUser(ctx context.Context) (*GetOrCreateDefaultUserResponse, error)
}

type GetUserByIDRequest struct {
	UserID string `json:"user_id"`
}

type GetUserByIDResponse struct {
	User *models.User `json:"user"`
}

type GetOrCreateDefaultUserResponse struct {
	User *models.User `json:"user"`
}
