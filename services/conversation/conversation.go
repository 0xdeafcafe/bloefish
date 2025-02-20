package conversation

import (
	"context"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
)

type Service interface {
	CreateConversation(ctx context.Context, req *CreateConversationRequest) (*CreateConversationResponse, error)
	CreateConversationMessage(ctx context.Context, req *CreateConversationMessageRequest) (*CreateConversationMessageResponse, error)
	GetInteraction(ctx context.Context, req *GetInteractionRequest) (*GetInteractionResponse, error)
	GetConversationWithInteractions(ctx context.Context, req *GetConversationWithInteractionsRequest) (*GetConversationWithInteractionsResponse, error)
	ListConversationsWithInteractions(ctx context.Context, req *ListConversationsWithInteractionsRequest) (*ListConversationsWithInteractionsResponse, error)
	DeleteConversations(ctx context.Context, req *DeleteConversationsRequest) error
	DeleteInteractions(ctx context.Context, req *DeleteInteractionsRequest) error
	UpdateInteractionExcludedState(ctx context.Context, req *UpdateInteractionExcludedStateRequest) error
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
	ID             string          `json:"id"`
	Owner          *Actor          `json:"owner"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`

	Title           *string `json:"title"`
	StreamChannelID string  `json:"stream_channel_id"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
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
	ConversationID      string                                        `json:"conversation_id"`
	InputInteraction    *CreateConversationMessageResponseInteraction `json:"input_interaction"`
	ResponseInteraction *CreateConversationMessageResponseInteraction `json:"response_interaction"`
	StreamChannelID     string                                        `json:"stream_channel_id"`
}

type CreateConversationMessageResponseInteraction struct {
	ID              string   `json:"id"`
	FileIDs         []string `json:"file_ids"`
	StreamChannelID string   `json:"stream_channel_id"`

	MarkedAsExcludedAt *time.Time `json:"marked_as_excluded_at"`

	MessageContent string   `json:"message_content"`
	Errors         []cher.E `json:"errors"`

	Owner          *Actor          `json:"owner"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type GetInteractionRequest struct {
	InteractionID string `json:"interaction_id"`
}

type GetInteractionResponse struct {
	ID             string   `json:"id"`
	ConversationID string   `json:"conversation_id"`
	FileIDs        []string `json:"file_ids"`

	MarkedAsExcludedAt *time.Time `json:"marked_as_excluded_at"`

	MessageContent string   `json:"message_content"`
	Errors         []cher.E `json:"errors"`

	Owner          *Actor          `json:"owner"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type GetConversationWithInteractionsRequest struct {
	ConversationID string `json:"conversation_id"`
}

type GetConversationWithInteractionsResponse struct {
	ConversationID string          `json:"conversation_id"`
	Owner          *Actor          `json:"owner"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`

	Title           *string `json:"title"`
	StreamChannelID string  `json:"stream_channel_id"`

	Interactions []*GetConversationWithInteractionsResponseInteraction `json:"interactions"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type GetConversationWithInteractionsResponseInteraction struct {
	ID              string   `json:"id"`
	FileIDs         []string `json:"file_ids"`
	StreamChannelID string   `json:"stream_channel_id"`

	MarkedAsExcludedAt *time.Time `json:"marked_as_excluded_at"`

	MessageContent string   `json:"message_content"`
	Errors         []cher.E `json:"errors"`

	Owner          *Actor          `json:"owner"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type ListConversationsWithInteractionsRequest struct {
	Owner *Actor `json:"owner"`
}

type ListConversationsWithInteractionsResponse struct {
	Conversations []*ListConversationsWithInteractionsResponseConversation `json:"conversations"`
}

type ListConversationsWithInteractionsResponseConversation struct {
	ID             string          `json:"id"`
	Owner          *Actor          `json:"owner"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`

	Title           *string `json:"title"`
	StreamChannelID string  `json:"stream_channel_id"`

	Interactions []*ListConversationsWithInteractionsResponseConversationInteraction `json:"interactions"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type ListConversationsWithInteractionsResponseConversationInteraction struct {
	ID              string   `json:"id"`
	FileIDs         []string `json:"file_ids"`
	StreamChannelID string   `json:"stream_channel_id"`

	MarkedAsExcludedAt *time.Time `json:"marked_as_excluded_at"`

	MessageContent string   `json:"message_content"`
	Errors         []cher.E `json:"errors"`

	Owner          *Actor          `json:"owner"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type DeleteConversationsRequest struct {
	ConversationIDs []string `json:"conversation_ids"`
}

type DeleteInteractionsRequest struct {
	InteractionIDs []string `json:"interaction_ids"`
}

type UpdateInteractionExcludedStateRequest struct {
	InteractionID string `json:"interaction_id"`
	Excluded      bool   `json:"excluded"`
}
