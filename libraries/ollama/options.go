package ollama

type ClientOption func(*client)

func WithEndpointURL(endpoint string) ClientOption {
	return func(c *client) {
		c.endpointURL = endpoint
	}
}
