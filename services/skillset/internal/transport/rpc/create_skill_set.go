package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/skillset"
)

func (r *RPC) CreateSkillSet(ctx context.Context, req *skillset.CreateSkillSetRequest) error {
	return r.app.CreateSkillSet(ctx, req)
}
