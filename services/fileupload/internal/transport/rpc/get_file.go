package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/fileupload"
)

func (r *RPC) GetFile(ctx context.Context, req *fileupload.GetFileRequest) (*fileupload.GetFileResponse, error) {
	return r.app.GetFile(ctx, req)
}
