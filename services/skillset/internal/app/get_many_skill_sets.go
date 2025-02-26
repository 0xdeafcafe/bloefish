package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/services/skillset"
)

func (a *App) GetManySkillSets(ctx context.Context, req *skillset.GetManySkillSetsRequest) (*skillset.GetManySkillSetsResponse, error) {
	skillSets, err := a.SkillSetRepository.GetManySkillSets(ctx, req.SkillSetIDs)
	if err != nil {
		return nil, err
	}

	if req.Owner != nil {
		for _, skillSet := range skillSets {
			if string(skillSet.Owner.Type) != string(req.Owner.Type) || skillSet.Owner.Identifier != req.Owner.Identifier {
				return nil, cher.New("skill_set_not_found", cher.M{"skill_set_id": skillSet.ID})
			}
		}
	}

	if !req.AllowDeleted {
		for _, skillSet := range skillSets {
			if skillSet.DeletedAt != nil {
				return nil, cher.New("skill_set_not_found", cher.M{"skill_set_id": skillSet.ID})
			}
		}
	}

	resp := &skillset.GetManySkillSetsResponse{
		SkillSets: make([]*skillset.SkillSet, len(skillSets)),
	}

	for i, skillSet := range skillSets {
		resp.SkillSets[i] = &skillset.SkillSet{
			ID:          skillSet.ID,
			Name:        skillSet.Name,
			Icon:        skillSet.Icon,
			Prompt:      skillSet.Prompt,
			Description: skillSet.Description,

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
