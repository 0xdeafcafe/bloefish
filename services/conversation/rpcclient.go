package conversation

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

func (r *RPCClient) CreateConversation(ctx context.Context, req *CreateConversationRequest) (resp *CreateConversationResponse, err error) {
	return resp, r.client.Do(ctx, "create_conversation", "2025-02-12", req, &resp)
}

func (r *RPCClient) CreateConversationMessage(ctx context.Context, req *CreateConversationMessageRequest) (resp *CreateConversationMessageResponse, err error) {
	return resp, r.client.Do(ctx, "create_conversation_message", "2025-02-12", req, &resp)
}

func (r *RPCClient) GetInteraction(ctx context.Context, req *GetInteractionRequest) (resp *GetInteractionResponse, err error) {
	return resp, r.client.Do(ctx, "get_interaction", "2025-02-12", req, &resp)
}

func (r *RPCClient) GetConversationWithInteractions(ctx context.Context, req *GetConversationWithInteractionsRequest) (resp *GetConversationWithInteractionsResponse, err error) {
	return resp, r.client.Do(ctx, "get_conversation_with_interactions", "2025-02-12", req, &resp)
}

func (r *RPCClient) ListConversationsWithInteractions(ctx context.Context, req *ListConversationsWithInteractionsRequest) (resp *ListConversationsWithInteractionsResponse, err error) {
	return resp, r.client.Do(ctx, "list_conversations_with_interactions", "2025-02-12", req, &resp)
}

func (r *RPCClient) DeleteConversations(ctx context.Context, req *DeleteConversationsRequest) error {
	return r.client.Do(ctx, "delete_conversations", "2025-02-12", req, nil)
}

func (r *RPCClient) DeleteInteractions(ctx context.Context, req *DeleteInteractionsRequest) error {
	return r.client.Do(ctx, "delete_interactions", "2025-02-12", req, nil)
}

func (r *RPCClient) UpdateInteractionExcludedState(ctx context.Context, req *UpdateInteractionExcludedStateRequest) error {
	return r.client.Do(ctx, "update_interaction_excluded_state", "2025-02-12", req, nil)
}
