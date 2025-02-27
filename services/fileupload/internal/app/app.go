package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/services/fileupload"
	"github.com/0xdeafcafe/bloefish/services/fileupload/internal/domain/models"
	"github.com/0xdeafcafe/bloefish/services/fileupload/internal/domain/ports"
)

type App struct {
	FileRepository    ports.FileRepository
	FileObjectService ports.FileObjectService
}

func (a *App) CreateUpload(ctx context.Context, req *fileupload.CreateUploadRequest) (*fileupload.CreateUploadResponse, error) {
	fileID, err := a.FileRepository.CreateUpload(ctx, &models.CreateUploadCommand{
		Name:     req.Name,
		Size:     req.Size,
		MIMEType: req.MIMEType,
		Owner: &models.CreateUploadCommandActor{
			Type:       models.ActorType(req.Owner.Type),
			Identifier: req.Owner.Identifier,
		},
	})
	if err != nil {
		return nil, err
	}

	presignedURL, err := a.FileObjectService.CreatePresignedUploadURL(ctx, fileID)
	if err != nil {
		return nil, err
	}

	return &fileupload.CreateUploadResponse{
		ID:        fileID,
		UploadURL: presignedURL,
	}, nil
}

func (a *App) ConfirmUpload(ctx context.Context, req *fileupload.ConfirmUploadRequest) error {
	file, err := a.FileRepository.Get(ctx, req.FileID)
	if err != nil {
		return err
	}

	if file.DeletedAt != nil {
		return cher.New("file_deleted", nil)
	}
	if file.ConfirmedAt != nil {
		return cher.New("file_already_confirmed", nil)
	}

	if err := a.FileObjectService.EnsureFileUploadValidity(ctx, file.ID, file.MIMEType, int(file.Size)); err != nil {
		return err
	}

	_, err = a.FileRepository.ConfirmUpload(ctx, req.FileID)
	if err != nil {
		return err
	}

	return nil
}
