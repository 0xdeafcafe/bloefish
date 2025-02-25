package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/skillset"
)

func (a *App) ListSkillSetsByOwner(ctx context.Context, req *skillset.ListSkillSetsByOwnerRequest) (*skillset.ListSkillSetsByOwnerResponse, error) {
	skillSets, err := a.SkillSetRepository.ListSkillSetsByOwner(ctx, string(req.Owner.Type), req.Owner.Identifier)
	if err != nil {
		return nil, err
	}

	resp := &skillset.ListSkillSetsByOwnerResponse{
		SkillSets: make([]*skillset.ListSkillSetsByOwnerResponseSkillSet, len(skillSets)),
	}

	for i, skillSet := range skillSets {
		resp.SkillSets[i] = &skillset.ListSkillSetsByOwnerResponseSkillSet{
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
		}
	}

	return resp, nil
}
