package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type NewStreamingChatParams struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

func (c *client) NewStreamingChat(ctx context.Context, params NewStreamingChatParams) (*StreamingChatIterator, error) {
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpointURL+"/api/chat", bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.streamClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, c.handleErrorResponse(resp)
	}

	return NewStreamingChatIterator(&Iterator[StreamingChatEvent]{
		scanner: bufio.NewScanner(resp.Body),
	}), nil
}
