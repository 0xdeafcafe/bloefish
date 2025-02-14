package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/fileupload"
)

func (r *RPC) HandleCreateUpload(ctx context.Context, req *fileupload.CreateUploadRequest) (*fileupload.CreateUploadResponse, error) {
	return r.app.CreateUpload(ctx, req)
}
