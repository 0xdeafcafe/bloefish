package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/fileupload"
)

func (r *RPC) ConfirmUpload(ctx context.Context, req *fileupload.ConfirmUploadRequest) error {
	return r.app.ConfirmUpload(ctx, req)
}
