package ollama

import (
	"context"
	"net/http"
	"time"
)

type client struct {
	baseURL      string
	httpClient   *http.Client
	streamClient *http.Client
}

type Client interface {
	NewStreamingChat(ctx context.Context, params NewStreamingChatParams) (*StreamingChatIterator, error)
}

func NewClient() Client {
	return &client{
		baseURL: "http://localhost:11434",
		httpClient: &http.Client{
			Timeout: time.Second * 2, // None as we stream
		},
		streamClient: &http.Client{
			Timeout: 0, // None as we stream
		},
	}
}
