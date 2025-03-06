package relay

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/langwatch"
)

type Client struct {
	providers map[ProviderID]Provider
	langwatch langwatch.Client
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		providers: make(map[ProviderID]Provider),
		langwatch: langwatch.NewNoopClient(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) With(providerID string) Provider {
	provider := c.providers[ProviderID(providerID)]
	if provider == nil {
		return newUnknownProvider()
	}

	// Wrap the provider with langwatch capabilities
	return newLangwatchProvider(provider, c.langwatch)
}

func (c *Client) ListAllModels(ctx context.Context) ([]Model, error) {
	var models []Model
	for _, provider := range c.providers {
		providerModels, err := provider.ListModels(ctx)
		if err != nil {
			return nil, err
		}

		models = append(models, providerModels...)
	}

	return models, nil
}
