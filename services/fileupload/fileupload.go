package fileupload

import (
	"context"
)

type Service interface {
	CreateUpload(ctx context.Context, req *CreateUploadRequest) (*CreateUploadResponse, error)
	ConfirmUpload(ctx context.Context, req *ConfirmUploadRequest) error
	GetFile(ctx context.Context, req *GetFileRequest) (*GetFileResponse, error)
	GetManyFiles(ctx context.Context, req *GetManyFilesRequest) (*GetManyFilesResponse, error)
}

type ActorType string

const (
	ActorTypeUser ActorType = "user"
)

type Actor struct {
	Type       ActorType `json:"type"`
	Identifier string    `json:"identifier"`
}

type File struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Size               int64   `json:"size"`
	MIMEType           string  `json:"mime_type"`
	Owner              *Actor  `json:"owner"`
	PresignedAccessURL *string `json:"presigned_access_url"`
}

type CreateUploadRequest struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	MIMEType string `json:"mime_type"`
	Owner    *Actor `json:"owner"`
}

type CreateUploadResponse struct {
	ID        string `json:"id"`
	UploadURL string `json:"upload_url"`
}

type ConfirmUploadRequest struct {
	FileID string `json:"file_id"`
}

type GetFileRequest struct {
	FileID                 string `json:"file_id"`
	IncludeAccessURL       bool   `json:"include_access_url"`
	AccessURLExpirySeconds *int   `json:"access_url_expiry_seconds"`
}

type GetFileResponse struct {
	File
}

type GetManyFilesRequest struct {
	FileIDs                []string `json:"file_ids"`
	Owner                  *Actor   `json:"owner"`
	AllowDeleted           bool     `json:"allow_deleted"`
	IncludeAccessURL       bool     `json:"include_access_url"`
	AccessURLExpirySeconds *int     `json:"access_url_expiry_seconds"`
}

type GetManyFilesResponse struct {
	Files []*File `json:"files"`
}
