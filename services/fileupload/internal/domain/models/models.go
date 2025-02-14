package models

import "time"

type File struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	MIMEType string `json:"mime_type"`
	Owner    *Actor `json:"owner"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	ConfirmedAt *time.Time `json:"confirmed_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type CreateUploadCommand struct {
	Name     string
	Size     int64
	MIMEType string
	Owner    *CreateUploadCommandActor
}

type CreateUploadCommandActor struct {
	Type       ActorType
	Identifier string
}
