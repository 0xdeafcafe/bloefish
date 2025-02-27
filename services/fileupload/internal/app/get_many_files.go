package app

import (
	"context"
	"runtime"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/services/fileupload"
	"golang.org/x/sync/errgroup"
)

func (a *App) GetManyFiles(ctx context.Context, req *fileupload.GetManyFilesRequest) (*fileupload.GetManyFilesResponse, error) {
	files, err := a.FileRepository.GetMany(ctx, req.FileIDs)
	if err != nil {
		return nil, err
	}

	if req.Owner != nil {
		for _, file := range files {
			if string(file.Owner.Type) != string(req.Owner.Type) || file.Owner.Identifier != req.Owner.Identifier {
				return nil, cher.New("file_not_found", cher.M{"file_id": file.ID})
			}
		}
	}

	if !req.AllowDeleted {
		for _, file := range files {
			if file.DeletedAt != nil {
				return nil, cher.New("file_not_found", cher.M{"file_id": file.ID})
			}
		}
	}

	resp := &fileupload.GetManyFilesResponse{
		Files: make([]*fileupload.File, len(files)),
	}

	errGroup, egCtx := errgroup.WithContext(ctx)
	errGroup.SetLimit(runtime.NumCPU() * 4)

	for i, file := range files {
		errGroup.Go(func() error {
			resp.Files[i] = &fileupload.File{
				ID:       file.ID,
				Name:     file.Name,
				Size:     file.Size,
				MIMEType: file.MIMEType,
				Owner: &fileupload.Actor{
					Type:       fileupload.ActorType(file.Owner.Type),
					Identifier: file.Owner.Identifier,
				},
			}

			if req.IncludeAccessURL {
				expiry := time.Minute * 15
				if req.AccessURLExpirySeconds != nil {
					expirySeconds := *req.AccessURLExpirySeconds
					expiry = time.Second * time.Duration(expirySeconds)
				}

				presignedURL, err := a.FileObjectService.CreatePresignedDownloadURL(egCtx, file.ID, expiry)
				if err != nil {
					return err
				}

				resp.Files[i].PresignedAccessURL = &presignedURL
			}

			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	return resp, nil
}
