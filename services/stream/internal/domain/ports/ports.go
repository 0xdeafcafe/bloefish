package ports

import (
	"context"

	"github.com/gorilla/websocket"
)

type MessageBroker interface {
	RegisterConnection(ctx context.Context, conn *websocket.Conn)
	SendMessageFull(ctx context.Context, channelID, messageContent string) error
	SendMessageFragment(ctx context.Context, channelID, messageContent string) error
}
