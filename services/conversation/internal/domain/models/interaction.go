package models

import (
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
)

type Interaction struct {
	ID             string   `json:"id"`
	IdempotencyKey string   `json:"idempotency_key"`
	ConversationID string   `json:"conversation_id"`
	FileIDs        []string `json:"file_ids"`
	SkillSetIDs    []string `json:"skill_set_ids"`

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

type CreateInteractionCommand struct {
	IdempotencyKey string
	ConversationID string
	FileIDs        []string
	SkillSetIDs    []string
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

type CreateActiveInteractionCommand struct {
	IdempotencyKey string
	ConversationID string
	FileIDs        []string
	SkillSetIDs    []string
	MessageContent string
	Owner          *CreateActiveInteractionCommandOwner
	AIRelayOptions *CreateActiveInteractionCommandAIRelayOptions
}

type CreateActiveInteractionCommandOwner struct {
	Type       ActorType
	Identifier string
}

type CreateActiveInteractionCommandAIRelayOptions struct {
	ProviderID string
	ModelID    string
}
