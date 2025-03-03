package ollama

import (
	ollamaClient "github.com/0xdeafcafe/bloefish/libraries/ollama"
	"github.com/0xdeafcafe/bloefish/services/airelay/internal/libraries/relay"
)

func NewProvider(
	ollamaClient ollamaClient.Client,
) relay.Provider {
	return &Provider{
		client: ollamaClient,
	}
}

type Provider struct {
	client ollamaClient.Client
}

func (p *Provider) GetMetadata() relay.ProviderMetadata {
	return relay.ProviderMetadata{
		ProviderID: relay.ProviderIdOllama,
		Name:       "Ollama",
	}
}
