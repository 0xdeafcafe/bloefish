package app

import (
	"context"
	"io"
	"net/http"
	"runtime"
	"sync"

	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/0xdeafcafe/bloefish/services/fileupload"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/sync/errgroup"
)

type downloadedFile struct {
	fileupload.File
	Content []byte
}

func (a *App) downloadFiles(ctx context.Context, owner *airelay.Actor, fileIDs []string) (map[string]*downloadedFile, error) {
	if len(fileIDs) == 0 {
		return map[string]*downloadedFile{}, nil
	}

	files, err := a.FileUploadService.GetManyFiles(ctx, &fileupload.GetManyFilesRequest{
		FileIDs: fileIDs,
		Owner: &fileupload.Actor{
			Type:       fileupload.ActorType(owner.Type),
			Identifier: owner.Identifier,
		},
		AllowDeleted:     false,
		IncludeAccessURL: true,
	})
	if err != nil {
		return nil, err
	}

	egGroup, egCtx := errgroup.WithContext(ctx)
	egGroup.SetLimit(runtime.NumCPU() * 4)

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	downloadedFiles := make(map[string]*downloadedFile, len(files.Files))
	mu := sync.Mutex{}

	for _, file := range files.Files {
		egGroup.Go(func() error {
			req, err := http.NewRequestWithContext(egCtx, http.MethodGet, *file.PresignedAccessURL, nil)
			if err != nil {
				return err
			}

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			content, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()
			downloadedFiles[file.ID] = &downloadedFile{
				File:    *file,
				Content: content,
			}

			return nil
		})
	}

	return nil, nil
}
