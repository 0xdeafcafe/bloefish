package relay

import (
	"github.com/0xdeafcafe/bloefish/libraries/ollama"
	"github.com/openai/openai-go"
)

type Client struct {
	ollamaClient ollama.Client
	openAIClient *openai.Client

	providers []Provider
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		ollamaClient: ollama.NewClient(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Providers() []Provider {
	return c.providers
}
