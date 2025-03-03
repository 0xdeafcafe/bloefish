package airelay

import "context"

type Service interface {
	ListSupported(ctx context.Context) (*ListSupportedResponse, error)
	InvokeConversationMessage(ctx context.Context, req *InvokeConversationMessageRequest) (*InvokeConversationMessageResponse, error)
	InvokeStreamingConversationMessage(ctx context.Context, req *InvokeStreamingConversationMessageRequest) (*InvokeStreamingConversationMessageResponse, error)
}

type ActorType string

const (
	ActorTypeUser ActorType = "user"
	ActorTypeBot  ActorType = "bot"
)

type Actor struct {
	Type       ActorType `json:"type"`
	Identifier string    `json:"identifier"`
}

type ListSupportedResponse struct {
	Models []*ListSupportedResponseModel `json:"models"`
}

type ListSupportedResponseModel struct {
	ProviderID   string `json:"provider_id"`
	ProviderName string `json:"provider_name"`
	ModelID      string `json:"model_id"`
	ModelName    string `json:"model_name"`
}

type InvokeConversationMessageRequest struct {
	Owner          *Actor                                          `json:"owner"`
	Messages       []*InvokeConversationMessageRequestMessage      `json:"messages"`
	AIRelayOptions *InvokeConversationMessageRequestAIRelayOptions `json:"ai_relay_options"`
}

type InvokeConversationMessageRequestMessage struct {
	Content string   `json:"content"`
	Owner   *Actor   `json:"owner"`
	FileIDs []string `json:"file_ids"`
}

type InvokeConversationMessageRequestAIRelayOptions struct {
	ProviderID string `json:"provider_id"`
	ModelID    string `json:"model_id"`
}

type InvokeConversationMessageResponse struct {
	MessageContent string `json:"message_content"`
}
type InvokeStreamingConversationMessageRequest struct {
	StreamingChannelID string                                          `json:"streaming_channel_id"`
	Owner              *Actor                                          `json:"owner"`
	Messages           []*InvokeConversationMessageRequestMessage      `json:"messages"`
	AIRelayOptions     *InvokeConversationMessageRequestAIRelayOptions `json:"ai_relay_options"`
}

type InvokeStreamingConversationMessageResponse struct {
	MessageContent string `json:"message_content"`
}
