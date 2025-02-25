package user

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

func (c *RPCClient) GetUserByID(ctx context.Context, req *GetUserByIDRequest) (resp *GetUserByIDResponse, err error) {
	return resp, c.client.Do(ctx, "get_user_by_id", "2025-02-12", req, &resp)
}

func (c *RPCClient) GetOrCreateDefaultUser(ctx context.Context) (resp *GetOrCreateDefaultUserResponse, err error) {
	return resp, c.client.Do(ctx, "get_or_create_default_user", "2025-02-12", nil, &resp)
}
