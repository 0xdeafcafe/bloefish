package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/skillset"
)

func (a *App) GetSkillSet(ctx context.Context, req *skillset.GetSkillSetRequest) (*skillset.GetSkillSetResponse, error) {
	skillSet, err := a.SkillSetRepository.GetSkillSet(ctx, req.SkillSetID)
	if err != nil {
		return nil, err
	}

	return &skillset.GetSkillSetResponse{
		SkillSet: skillset.SkillSet{
			ID:          skillSet.ID,
			Name:        skillSet.Name,
			Icon:        skillSet.Icon,
			Description: skillSet.Description,
			Prompt:      skillSet.Prompt,

			Owner: &skillset.Actor{
				Type:       skillset.ActorType(skillSet.Owner.Type),
				Identifier: skillSet.Owner.Identifier,
			},

			CreatedAt: skillSet.CreatedAt,
			UpdatedAt: skillSet.UpdatedAt,
			DeletedAt: skillSet.DeletedAt,
		},
	}, nil
}
