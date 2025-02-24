package relay

import (
	"context"
	"errors"
	"net"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/ollama"
)

type ollamaChatStreamIterator struct {
	inner *ollama.StreamingChatIterator
}

func (c *Client) newOllamaChatStream(ctx context.Context, params ChatStreamParams) (ChatStreamIterator, error) {
	messages := make([]ollama.Message, len(params.Messages))
	for i, msg := range params.Messages {
		messages[i] = ollama.Message{
			Role:    ollama.Role(msg.Role),
			Content: msg.Content,
		}
	}

	stream, err := c.ollamaClient.NewStreamingChat(ctx, ollama.NewStreamingChatParams{
		Model:    params.ModelID,
		Messages: messages,
	})
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			return nil, cher.New("ollama_connection_issue", nil, cher.Coerce(err))
		}

		return nil, err
	}

	return &ollamaChatStreamIterator{
		inner: stream,
	}, nil
}

func (i *ollamaChatStreamIterator) Next() bool {
	return i.inner.Next()
}

func (i *ollamaChatStreamIterator) Current() *ChatStreamEvent {
	current := i.inner.Current()
	return &ChatStreamEvent{
		Content: current.Message.Content,
		Done:    current.Done,
	}
}

func (i *ollamaChatStreamIterator) Content() string {
	return i.inner.Content()
}

func (i *ollamaChatStreamIterator) Err() error {
	return i.inner.Err()
}
