package models

import "time"

type Interaction struct {
	ID             string   `json:"id"`
	IdempotencyKey string   `json:"idempotency_key"`
	ConversationID string   `json:"conversation_id"`
	MessageContent string   `json:"message_content"`
	FileIDs        []string `json:"file_ids"`

	Owner          *Actor          `json:"owner"`
	AIRelayOptions *AIRelayOptions `json:"ai_relay_options"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateInteractionCommand struct {
	IdempotencyKey string
	ConversationID string
	FileIDs        []string
	MessageContent string
	Owner          *CreateInteractionCommandOwner
	AIRelayOptions *CreateInteractionCommandAIRelayOptions
}

type CreateInteractionCommandOwner struct {
	Type       ActorType
	Identifier string
}

type CreateInteractionCommandAIRelayOptions struct {
	ProviderID string
	ModelID    string
}
