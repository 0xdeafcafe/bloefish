package openai

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/airelay/internal/libraries/relay"
)

type Model struct {
	ID   string
	Name string
}

func (p *Provider) ListModels(ctx context.Context) ([]relay.Model, error) {
	var result []relay.Model
	for _, model := range p.models {
		result = append(result, relay.Model{
			ProviderID: p.GetMetadata().ProviderID,
			ModelID:    model.ID,
			ModelName:  model.Name,
		})
	}

	return result, nil
}
