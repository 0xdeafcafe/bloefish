package openai

import (
	"github.com/0xdeafcafe/bloefish/services/airelay/internal/libraries/relay"
	oaiClient "github.com/openai/openai-go"
)

type Provider struct {
	client oaiClient.Client
	models []Model
}

func NewProvider(
	oaiClient oaiClient.Client,
	opts ...ProviderOption,
) relay.Provider {
	p := &Provider{
		client: oaiClient,
		models: []Model{},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Provider) GetMetadata() relay.ProviderMetadata {
	return relay.ProviderMetadata{
		ProviderID: relay.ProviderIdOpenAI,
		Name:       "Open AI",
	}
}
