package openai

type ProviderOption func(*Provider)

func WithModels(models []Model) ProviderOption {
	return func(c *Provider) {
		c.models = models
	}
}
