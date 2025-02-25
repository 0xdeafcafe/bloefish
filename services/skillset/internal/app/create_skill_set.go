package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/skillset"
	"github.com/0xdeafcafe/bloefish/services/skillset/internal/domain/models"
)

func (a *App) CreateSkillSet(ctx context.Context, req *skillset.CreateSkillSetRequest) error {
	if _, err := a.SkillSetRepository.CreateSkillSet(ctx, models.CreateSkillSetCommand{
		Name:        req.Name,
		Icon:        req.Icon,
		Description: req.Description,
		Prompt:      req.Prompt,
		Owner: &models.CreateSkillSetCommandActor{
			Type:       models.ActorType(req.Owner.Type),
			Identifier: req.Owner.Identifier,
		},
	}); err != nil {
		return err
	}

	return nil
}
