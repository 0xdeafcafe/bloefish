package airelay

import (
	"context"
	"net/http"

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

func (r *RPCClient) ListSupported(ctx context.Context) (resp *ListSupportedResponse, err error) {
	return resp, r.client.Do(ctx, "list_supported", "2025-02-12", nil, &resp)
}

func (r *RPCClient) InvokeConversationMessage(ctx context.Context, req *InvokeConversationMessageRequest) (resp *InvokeConversationMessageResponse, err error) {
	return resp, r.client.Do(ctx, "invoke_conversation_message", "2025-02-12", req, &resp)
}

func (r *RPCClient) InvokeStreamingConversationMessage(ctx context.Context, req *InvokeStreamingConversationMessageRequest) (resp *InvokeStreamingConversationMessageResponse, err error) {
	return resp, r.client.Do(ctx, "invoke_streaming_conversation_message", "2025-02-12", req, &resp)
}
