package models

import "time"

type SkillSet struct {
	ID          string
	Name        string
	Icon        string
	Description string
	Prompt      string

	Owner *Actor

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type CreateSkillSetCommand struct {
	Name        string
	Icon        string
	Description string
	Prompt      string
	Owner       *CreateSkillSetCommandActor
}

type CreateSkillSetCommandActor struct {
	Type       ActorType
	Identifier string
}
