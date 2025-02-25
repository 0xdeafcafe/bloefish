package skillset

import (
	"context"
	"net/http"

	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/libraries/crpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type RPCClient struct {
	client *crpc.Client
}

func NewRPCClient(ctx context.Context, cfg config.UnauthenticatedService) Service {
	return &RPCClient{
		client: crpc.NewClient(ctx, cfg.BaseURL, &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}),
	}
}

func (r *RPCClient) CreateSkillSet(ctx context.Context, req *CreateSkillSetRequest) error {
	return r.client.Do(ctx, "create_skill_set", "2025-02-12", req, nil)
}

func (r *RPCClient) GetSkillSet(ctx context.Context, req *GetSkillSetRequest) (resp *GetSkillSetResponse, err error) {
	return resp, r.client.Do(ctx, "get_skill_set", "2025-02-12", req, &resp)
}

func (r *RPCClient) ListSkillSetsByOwner(ctx context.Context, req *ListSkillSetsByOwnerRequest) (resp *ListSkillSetsByOwnerResponse, err error) {
	return resp, r.client.Do(ctx, "list_skill_sets_by_owner", "2025-02-12", req, &resp)
}
