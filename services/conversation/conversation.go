package conversation

import (
	"context"
	"time"
)

type Service interface {
	CreateConversation(ctx context.Context, req *CreateConversationRequest) (*CreateConversationResponse, error)
	CreateConversationMessage(ctx context.Context, req *CreateConversationMessageRequest) (*CreateConversationMessageResponse, error)
	GetInteraction(ctx context.Context, req *GetInteractionRequest) (*GetInteractionResponse, error)
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

type AIRelayOptions struct {
	ProviderID string `json:"provider_id"`
	ModelID    string `json:"model_id"`
}

type CreateConversationRequest struct {
	IdempotencyKey string          `json:"idempotency_key"`
	Owner          *Actor          `json:"owner"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`
}

type CreateConversationResponse struct {
	ConversationID  string `json:"conversation_id"`
	StreamChannelID string `json:"stream_channel_id"`
}

type CreateConversationMessageRequest struct {
	ConversationID string                                   `json:"conversation_id"`
	IdempotencyKey string                                   `json:"idempotency_key"`
	MessageContent string                                   `json:"message_content"`
	FileIDs        []string                                 `json:"file_ids"`
	Owner          *Actor                                   `json:"owner"`
	AIRelayOptions *AIRelayOptions                          `json:"ai_relay_options"`
	Options        *CreateConversationMessageRequestOptions `json:"options"`
}

type CreateConversationMessageRequestOptions struct {
	UseStreaming bool `json:"use_streaming"`
}

type CreateConversationMessageResponse struct {
	ConversationID  string `json:"conversation_id"`
	InteractionID   string `json:"interaction_id"`
	StreamChannelID string `json:"stream_channel_id"`
}

type GetInteractionRequest struct {
	InteractionID string `json:"interaction_id"`
}

type GetInteractionResponse struct {
	ID             string          `json:"id"`
	ConversationID string          `json:"conversation_id"`
	Owner          *Actor          `json:"owner"`
	MessageContent string          `json:"message_content"`
	FileIDs        []string        `json:"file_ids"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      *time.Time      `json:"updated_at"`
}
