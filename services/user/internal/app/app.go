package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/user/internal/domain/models"
	"github.com/0xdeafcafe/bloefish/services/user/internal/domain/ports"
)

type App struct {
	UserRepository ports.UserRepository
}

func (a *App) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return a.UserRepository.GetByUserID(ctx, userID)
}

func (a *App) GetOrCreateDefaultUser(ctx context.Context) (*models.User, error) {
	return a.UserRepository.GetOrCreateDefaultUser(ctx)
}
