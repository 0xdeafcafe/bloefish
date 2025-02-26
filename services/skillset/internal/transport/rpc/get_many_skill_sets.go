package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/skillset"
)

func (r *RPC) GetManySkillSets(ctx context.Context, req *skillset.GetManySkillSetsRequest) (*skillset.GetManySkillSetsResponse, error) {
	return r.app.GetManySkillSets(ctx, req)
}
