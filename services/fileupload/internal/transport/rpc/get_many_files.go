package rpc

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/fileupload"
)

func (r *RPC) GetManyFiles(ctx context.Context, req *fileupload.GetManyFilesRequest) (*fileupload.GetManyFilesResponse, error) {
	return r.app.GetManyFiles(ctx, req)
}
