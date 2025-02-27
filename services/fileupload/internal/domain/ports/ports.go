package ports

import (
	"context"
	"time"

	"github.com/0xdeafcafe/bloefish/services/fileupload/internal/domain/models"
)

type FileRepository interface {
	CreateUpload(ctx context.Context, req *models.CreateUploadCommand) (string, error)
	ConfirmUpload(ctx context.Context, fileID string) (*models.File, error)
	Get(ctx context.Context, fileID string) (*models.File, error)
	GetMany(ctx context.Context, ids []string) ([]*models.File, error)
}

type FileObjectService interface {
	CreatePresignedUploadURL(ctx context.Context, fileID string) (string, error)
	CreatePresignedDownloadURL(ctx context.Context, fileID string, expiry time.Duration) (string, error)

	EnsureFileUploadValidity(ctx context.Context, fileID, mimeType string, size int) error
}
