package relay

import (
	"context"
	"errors"
)

type ProviderID string

const (
	ProviderIdOpenAI  ProviderID = "open_ai"
	ProviderIdOllama  ProviderID = "ollama"
	providerIdUnknown ProviderID = "unknown"
)

var (
	ErrRequiredProviderMissing = errors.New("required provider is missing")
)

type Provider interface {
	NewChatStream(ctx context.Context, params ChatStreamParams) (iter ChatStreamIterator, err error)
	ListModels(ctx context.Context) ([]Model, error)
	GetMetadata() ProviderMetadata
}

type ProviderMetadata struct {
	ProviderID ProviderID
	Name       string
}

func newUnknownProvider() Provider {
	return &unknownProvider{}
}

type unknownProvider struct{}

func (p *unknownProvider) NewChatStream(context.Context, ChatStreamParams) (ChatStreamIterator, error) {
	return nil, ErrRequiredProviderMissing
}

func (p *unknownProvider) ListModels(context.Context) ([]Model, error) {
	return nil, ErrRequiredProviderMissing
}

func (p *unknownProvider) GetMetadata() ProviderMetadata {
	return ProviderMetadata{
		ProviderID: providerIdUnknown,
		Name:       "Unknown",
	}
}
