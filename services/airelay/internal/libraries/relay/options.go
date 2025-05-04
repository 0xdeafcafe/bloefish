package relay

type ClientOption func(*Client)

func WithProvider(provider Provider) ClientOption {
	return func(c *Client) {
		c.providers[provider.GetMetadata().ProviderID] = provider
	}
}
