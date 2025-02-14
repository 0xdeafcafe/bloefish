package models

import "time"

type Conversation struct {
	ID             string                      `json:"id"`
	Owner          *Actor                      `json:"owner"`
	IdempotencyKey string                      `json:"idempotency_key"`
	AIRelayOptions *ConversationAIRelayOptions `json:"ai_relay_options"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type ConversationAIRelayOptions struct {
	ProviderID string `json:"provider_id"`
	ModelID    string `json:"model_id"`
}

type CreateConversationCommand struct {
	IdempotencyKey string
	Owner          *CreateConversationCommandOwner
	AIRelayOptions *CreateConversationCommandAIRelayOptions
}

type CreateConversationCommandOwner struct {
	Type       ActorType
	Identifier string
}

type CreateConversationCommandAIRelayOptions struct {
	ProviderID string
	ModelID    string
}
