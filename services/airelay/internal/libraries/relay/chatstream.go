package relay

import (
	"context"
	"fmt"
)

type ChatStreamParams struct {
	ProviderID string
	ModelID    string
	Messages   []Message
}

type ChatStreamEvent struct {
	Content string
	Done    bool
}

type ChatStreamIterator interface {
	Next() bool
	Current() *ChatStreamEvent
	Content() string
	Err() error
}

func (c *Client) NewChatStream(ctx context.Context, params ChatStreamParams) (ChatStreamIterator, error) {
	switch params.ProviderID {
	case "ollama":
		return c.newOllamaChatStream(ctx, params)
	case "open_ai":
		return c.newOpenAIChatStream(ctx, params)

	default:
		return nil, fmt.Errorf("%v: %s", ErrUnsupportedProvider, params.ProviderID)
	}
}
