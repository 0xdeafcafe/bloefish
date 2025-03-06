package langwatch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	defaultEndpointURL                   = "https://app.langwatch.ai"
	defaultAPIKeyEnvironmentVariableName = "LANGWATCH_API_KEY"
)

type client struct {
	endpointURL string
	apiKey      string
	httpClient  *http.Client
}

type Client interface {
	CreateTrace(ctx context.Context, threadID string, opts ...TraceOption) (context.Context, Trace)
}

func NewClient(opts ...ClientOption) Client {
	c := &client{
		endpointURL: defaultEndpointURL,
		apiKey:      os.Getenv(defaultAPIKeyEnvironmentVariableName),
		httpClient:  &http.Client{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *client) Collect(ctx context.Context, body interface{}) error {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal collector request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpointURL+"/api/collector", bytes.NewReader(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to create collector request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make collector request: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to successfully make collector request: %w", err)
	}

	return nil
}
