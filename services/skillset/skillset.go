package skillset

import (
	"context"
	"time"
)

type Service interface {
	CreateSkillSet(ctx context.Context, req *CreateSkillSetRequest) error
	GetSkillSet(ctx context.Context, req *GetSkillSetRequest) (*GetSkillSetResponse, error)
	GetManySkillSets(ctx context.Context, req *GetManySkillSetsRequest) (*GetManySkillSetsResponse, error)
	ListSkillSetsByOwner(ctx context.Context, req *ListSkillSetsByOwnerRequest) (*ListSkillSetsByOwnerResponse, error)
}

type ActorType string

const (
	ActorTypeUser ActorType = "user"
)

type SkillSet struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`

	Owner *Actor `json:"owner"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type Actor struct {
	Type       ActorType `json:"type"`
	Identifier string    `json:"identifier"`
}

type CreateSkillSetRequest struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
	Owner       *Actor `json:"owner"`
}

type GetSkillSetRequest struct {
	SkillSetID string `json:"skill_set_id"`
}

type GetSkillSetResponse struct {
	SkillSet
}

type GetManySkillSetsRequest struct {
	SkillSetIDs  []string `json:"skill_set_ids"`
	Owner        *Actor   `json:"owner"`
	AllowDeleted bool     `json:"allow_deleted"`
}

type GetManySkillSetsResponse struct {
	SkillSets []*SkillSet `json:"skill_sets"`
}

type ListSkillSetsByOwnerRequest struct {
	Owner *Actor `json:"owner"`
}

type ListSkillSetsByOwnerResponse struct {
	SkillSets []*SkillSet `json:"skill_sets"`
}
