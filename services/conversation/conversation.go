package conversation

import (
	"context"
	"time"
)

type Service interface {
	CreateConversation(ctx context.Context, req *CreateConversationRequest) (*CreateConversationResponse, error)
	CreateConversationMessage(ctx context.Context, req *CreateConversationMessageRequest) (*CreateConversationMessageResponse, error)
	GetInteraction(ctx context.Context, req *GetInteractionRequest) (*GetInteractionResponse, error)
	GetConversationWithInteractions(ctx context.Context, req *GetConversationWithInteractionsRequest) (*GetConversationWithInteractionsResponse, error)
	ListConversationsWithInteractions(ctx context.Context, req *ListConversationsWithInteractionsRequest) (*ListConversationsWithInteractionsResponse, error)
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
	ConversationID        string `json:"conversation_id"`
	StreamChannelIDPrefix string `json:"stream_channel_id_prefix"`
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
	ConversationID        string `json:"conversation_id"`
	InteractionID         string `json:"interaction_id"`
	ResponseInteractionID string `json:"response_interaction_id"`
	StreamChannelID       string `json:"stream_channel_id"`
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

type GetConversationWithInteractionsRequest struct {
	ConversationID string `json:"conversation_id"`
}

type GetConversationWithInteractionsResponse struct {
	ConversationID string          `json:"conversation_id"`
	Owner          *Actor          `json:"owner"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`

	Interactions []*GetConversationWithInteractionsResponseInteraction `json:"interactions"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type GetConversationWithInteractionsResponseInteraction struct {
	ID             string          `json:"id"`
	Owner          *Actor          `json:"owner"`
	MessageContent string          `json:"message_content"`
	FileIDs        []string        `json:"file_ids"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      *time.Time      `json:"updated_at"`
	CompletedAt    *time.Time      `json:"completed_at"`
	DeletedAt      *time.Time      `json:"deleted_at"`
}

type ListConversationsWithInteractionsRequest struct {
	Owner *Actor `json:"owner"`
}

type ListConversationsWithInteractionsResponse struct {
	Conversations []*ListConversationsWithInteractionsResponseConversation `json:"conversations"`
}

type ListConversationsWithInteractionsResponseConversation struct {
	ID             string                                                              `json:"id"`
	Owner          *Actor                                                              `json:"owner"`
	AIRelayOptions *AIRelayOptions                                                     `json:"ai_relay_options"`
	Interactions   []*ListConversationsWithInteractionsResponseConversationInteraction `json:"interactions"`
	CreatedAt      time.Time                                                           `json:"created_at"`
	DeletedAt      *time.Time                                                          `json:"deleted_at"`
}

type ListConversationsWithInteractionsResponseConversationInteraction struct {
	ID             string          `json:"id"`
	Owner          *Actor          `json:"owner"`
	MessageContent string          `json:"message_content"`
	FileIDs        []string        `json:"file_ids"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      *time.Time      `json:"updated_at"`
	CompletedAt    *time.Time      `json:"completed_at"`
	DeletedAt      *time.Time      `json:"deleted_at"`
}
