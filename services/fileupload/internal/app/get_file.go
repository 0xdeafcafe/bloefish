package app

import (
	"context"
	"time"

	"github.com/0xdeafcafe/bloefish/services/fileupload"
)

func (a *App) GetFile(ctx context.Context, req *fileupload.GetFileRequest) (*fileupload.GetFileResponse, error) {
	file, err := a.FileRepository.Get(ctx, req.FileID)
	if err != nil {
		return nil, err
	}

	var presignedAccessURL *string
	if req.IncludeAccessURL {
		expiry := time.Minute * 15
		if req.AccessURLExpirySeconds != nil {
			expirySeconds := *req.AccessURLExpirySeconds
			expiry = time.Second * time.Duration(expirySeconds)
		}

		presignedURL, err := a.FileObjectService.CreatePresignedDownloadURL(ctx, file.ID, expiry)
		if err != nil {
			return nil, err
		}

		presignedAccessURL = &presignedURL
	}

	return &fileupload.GetFileResponse{
		File: fileupload.File{
			ID:       file.ID,
			Name:     file.Name,
			Size:     file.Size,
			MIMEType: file.MIMEType,
			Owner: &fileupload.Actor{
				Type:       fileupload.ActorType(file.Owner.Type),
				Identifier: file.Owner.Identifier,
			},
			PresignedAccessURL: presignedAccessURL,
		},
	}, nil
}
