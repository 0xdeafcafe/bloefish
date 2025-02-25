package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/skillset"
)

func (r *RPC) ListSkillSetsByOwner(ctx context.Context, req *skillset.ListSkillSetsByOwnerRequest) (*skillset.ListSkillSetsByOwnerResponse, error) {
	return r.app.ListSkillSetsByOwner(ctx, req)
}
