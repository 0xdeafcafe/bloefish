package services

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/ksuid"
	"github.com/0xdeafcafe/bloefish/libraries/merr"
	"github.com/0xdeafcafe/bloefish/services/stream/internal/domain/models"
	"github.com/0xdeafcafe/bloefish/services/stream/internal/domain/ports"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

type websocketMessageBroker struct {
	connections map[string]*websocket.Conn

	connectionsMu sync.Mutex
	writeMu       sync.Mutex
}

func NewWebSocketMessageBroker() ports.MessageBroker {
	return &websocketMessageBroker{
		connections: make(map[string]*websocket.Conn),

		connectionsMu: sync.Mutex{},
		writeMu:       sync.Mutex{},
	}
}

func (w *websocketMessageBroker) RegisterConnection(ctx context.Context, conn *websocket.Conn) {
	w.connectionsMu.Lock()
	defer w.connectionsMu.Unlock()
	w.connections[ksuid.Generate(ctx, "wsconn").String()] = conn
}

func (w *websocketMessageBroker) SendMessageFragment(ctx context.Context, channelID string, messageContent string) error {
	return w.sendMessage(ctx, channelID, &messageContent, nil, models.StreamMessageTypeMessageFragment)
}

func (w *websocketMessageBroker) SendMessageFull(ctx context.Context, channelID string, messageContent string) error {
	return w.sendMessage(ctx, channelID, &messageContent, nil, models.StreamMessageTypeMessageFull)
}

func (w *websocketMessageBroker) SendErrorMessage(ctx context.Context, channelID string, errorMessage cher.E) error {
	return w.sendMessage(ctx, channelID, nil, &errorMessage, models.StreamMessageTypeError)
}

func (w *websocketMessageBroker) sendMessage(ctx context.Context, channelID string, messageContent *string, errorMessage *cher.E, messageType models.StreamMessageType) error {
	msg := &models.StreamMessage{
		ChannelID: channelID,
		Type:      messageType,
	}

	switch msg.Type {
	case models.StreamMessageTypeMessageFragment:
		msg.MessageFragment = messageContent
	case models.StreamMessageTypeMessageFull:
		msg.MessageFull = messageContent
	case models.StreamMessageTypeError:
		msg.Error = errorMessage
	default:
		return merr.New(ctx, "invalid message type", merr.M{
			"message_type": messageType,
		})
	}

	jsonText, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal websocket message: %w", err)
	}

	errGroup, egCtx := errgroup.WithContext(ctx)
	errGroup.SetLimit(runtime.NumCPU())

	for connID, conn := range w.connections {
		errGroup.Go(func() error {
			w.writeMu.Lock()
			defer w.writeMu.Unlock()

			if err := conn.WriteMessage(websocket.TextMessage, jsonText); err != nil {
				w.connectionsMu.Lock()
				defer w.connectionsMu.Unlock()

				delete(w.connections, connID)

				clog.Get(egCtx).WithError(err).Warn("failed to write message to websocket connection")
			}

			return nil
		})
	}

	return nil
}
