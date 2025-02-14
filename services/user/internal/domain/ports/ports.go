package ports

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/user/internal/domain/models"
)

type UserRepository interface {
	GetByUserID(ctx context.Context, userID string) (*models.User, error)
	GetOrCreateDefaultUser(ctx context.Context) (*models.User, error)
}
