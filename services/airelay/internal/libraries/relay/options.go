package relay

import "github.com/0xdeafcafe/bloefish/libraries/langwatch"

type ClientOption func(*Client)

func WithProvider(provider Provider) ClientOption {
	return func(c *Client) {
		c.providers[provider.GetMetadata().ProviderID] = provider
	}
}

func WithLangwatch(langwatch langwatch.Client) ClientOption {
	return func(c *Client) {
		c.langwatch = langwatch
	}
}
