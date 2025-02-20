package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type NewStreamingChatParams struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

func (c *client) NewStreamingChat(ctx context.Context, params NewStreamingChatParams) (*StreamingChatIterator, error) {
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/chat", bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	resp, err := c.streamClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()

		return nil, err
	}

	return NewStreamingChatIterator(&Iterator[StreamingChatEvent]{
		scanner: bufio.NewScanner(resp.Body),
	}), nil
}
