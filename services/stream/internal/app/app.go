package app

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/stream"
	"github.com/0xdeafcafe/bloefish/services/stream/internal/domain/ports"
)

type App struct {
	MessageBroker ports.MessageBroker
}

func (a *App) SendMessageFull(ctx context.Context, req *stream.SendMessageFullRequest) error {
	return a.MessageBroker.SendMessageFull(ctx, req.ChannelID, req.MessageContent)
}

func (a *App) SendMessageFragment(ctx context.Context, req *stream.SendMessageFragmentRequest) error {
	return a.MessageBroker.SendMessageFragment(ctx, req.ChannelID, req.MessageContent)
}

func (a *App) SendErrorMessage(ctx context.Context, req *stream.SendErrorMessageRequest) error {
	return a.MessageBroker.SendErrorMessage(ctx, req.ChannelID, req.Error)
}
