package models

import "time"

type User struct {
	ID string `json:"id"`

	DefaultUser bool `json:"default_user"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
