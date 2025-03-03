package ollama

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/airelay/internal/libraries/relay"
)

func (p *Provider) ListModels(ctx context.Context) ([]relay.Model, error) {
	models, err := p.client.ListLocalModels(ctx)
	if err != nil {
		return nil, err
	}

	var result []relay.Model
	for _, model := range models {
		result = append(result, relay.Model{
			ProviderID: p.GetMetadata().ProviderID,
			ModelID:    model.Name,
			ModelName:  model.Name,
		})
	}

	return result, nil
}
