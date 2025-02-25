package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/skillset"
)

func (r *RPC) GetSkillSet(ctx context.Context, req *skillset.GetSkillSetRequest) (*skillset.GetSkillSetResponse, error) {
	return r.app.GetSkillSet(ctx, req)
}
