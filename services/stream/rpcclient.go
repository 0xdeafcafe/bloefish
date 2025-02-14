package stream

import (
	"context"
	"net/http"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/libraries/crpc"
)

type RPCClient struct {
	client *crpc.Client
}

func NewRPCClient(ctx context.Context, cfg config.UnauthenticatedService) Service {
	return &RPCClient{
		client: crpc.NewClient(ctx, cfg.BaseURL, &http.Client{}),
	}
}

func (r *RPCClient) SendMessageFragment(ctx context.Context, req *SendMessageFragmentRequest) error {
	return r.client.Do(ctx, "send_message_fragment", "2025-02-12", req, nil)
}

func (r *RPCClient) SendMessageFull(ctx context.Context, req *SendMessageFullRequest) error {
	return r.client.Do(ctx, "send_message_full", "2025-02-12", req, nil)
}

func (r *RPCClient) SendStreamedMessage(context.Context, *StreamedMessage) error {
	return cher.New("not_implemented", cher.M{"method": "SendStreamedMessage"})
}
