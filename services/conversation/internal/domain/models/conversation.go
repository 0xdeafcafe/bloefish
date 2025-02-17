package models

import "time"

type Conversation struct {
	ID             string          `json:"id"`
	Owner          *Actor          `json:"owner"`
	IdempotencyKey string          `json:"idempotency_key"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
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
