package relay

import (
	"github.com/0xdeafcafe/bloefish/libraries/ollama"
	"github.com/openai/openai-go"
)

type ClientOption func(*Client)

func WithOllamaClient(client ollama.Client) ClientOption {
	return func(c *Client) {
		c.ollamaClient = client
	}
}

func WithOpenAIClient(client *openai.Client) ClientOption {
	return func(c *Client) {
		c.openAIClient = client
	}
}

func WithProviderModels(providers []Provider) ClientOption {
	return func(c *Client) {
		c.providers = providers
	}
}
