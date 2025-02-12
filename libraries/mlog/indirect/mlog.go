package indirect

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/merr"
)

var (
	Debug func(ctx context.Context, err merr.Merrer)
	Info  func(ctx context.Context, err merr.Merrer)
	Warn  func(ctx context.Context, err merr.Merrer)
	Error func(ctx context.Context, err merr.Merrer)
)
