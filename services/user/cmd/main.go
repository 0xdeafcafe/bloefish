package cmd

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/user/internal"
)

func Root(ctx context.Context, args []string) error {
	return internal.Run(ctx)
}
