package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type client struct {
	endpointURL  string
	httpClient   *http.Client
	streamClient *http.Client
}

type Client interface {
	NewStreamingChat(ctx context.Context, params NewStreamingChatParams) (*StreamingChatIterator, error)
	ListLocalModels(ctx context.Context) ([]*LocalModel, error)
}

func NewClient(opts ...ClientOption) Client {
	c := &client{
		httpClient: &http.Client{
			Timeout: time.Second * 2, // None as we stream
		},
		streamClient: &http.Client{
			Timeout: 0, // None as we stream
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (c *client) handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read error response: %w", err)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("failed to unmarshal error response: %w", err)
	}

	return fmt.Errorf("ollama API error: %s", errResp.Error)
}
