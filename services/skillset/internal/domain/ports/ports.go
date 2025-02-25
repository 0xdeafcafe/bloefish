package ports

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/skillset/internal/domain/models"
)

type SkillSetRepository interface {
	CreateSkillSet(ctx context.Context, cmd models.CreateSkillSetCommand) (string, error)
	GetSkillSet(ctx context.Context, id string) (*models.SkillSet, error)
	ListSkillSetsByOwner(ctx context.Context, ownerType, ownerIdentifier string) ([]*models.SkillSet, error)
}
