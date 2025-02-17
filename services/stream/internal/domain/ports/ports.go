package ports

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/gorilla/websocket"
)

type MessageBroker interface {
	RegisterConnection(ctx context.Context, conn *websocket.Conn)
	SendMessageFull(ctx context.Context, channelID, messageContent string) error
	SendMessageFragment(ctx context.Context, channelID, messageContent string) error
	SendErrorMessage(ctx context.Context, channelID string, err cher.E) error
}
