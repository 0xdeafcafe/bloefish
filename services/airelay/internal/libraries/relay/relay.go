package relay

import (
	"context"
)

type Client struct {
	providers map[ProviderID]Provider
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		providers: make(map[ProviderID]Provider),
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

	return provider
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
