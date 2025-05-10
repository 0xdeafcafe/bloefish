package fileupload

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/libraries/crpc"
)

type RPCClient struct {
	client *crpc.Client
}

func NewRPCClient(ctx context.Context, cfg config.UnauthenticatedService) Service {
	return &RPCClient{
		client: crpc.NewClient(ctx, cfg.BaseURL, nil),
	}
}

func (r *RPCClient) CreateUpload(ctx context.Context, req *CreateUploadRequest) (resp *CreateUploadResponse, err error) {
	return resp, r.client.Do(ctx, "create_upload", "2025-02-12", req, &resp)
}

func (r *RPCClient) ConfirmUpload(ctx context.Context, req *ConfirmUploadRequest) error {
	return r.client.Do(ctx, "confirm_upload", "2025-02-12", req, nil)
}

func (r *RPCClient) GetFile(ctx context.Context, req *GetFileRequest) (resp *GetFileResponse, err error) {
	return resp, r.client.Do(ctx, "get_file", "2025-02-12", req, &resp)
}

func (r *RPCClient) GetManyFiles(ctx context.Context, req *GetManyFilesRequest) (resp *GetManyFilesResponse, err error) {
	return resp, r.client.Do(ctx, "get_many_files", "2025-02-12", req, &resp)
}
