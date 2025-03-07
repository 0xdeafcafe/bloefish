package relay

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/langwatch"
)

// langwatchProviderDecorator wraps a Provider and adds langwatch capabilities
type langwatchProviderDecorator struct {
	inner     Provider
	langwatch langwatch.Client
}

func newLangwatchProvider(provider Provider, langwatch langwatch.Client) Provider {
	return &langwatchProviderDecorator{
		inner:     provider,
		langwatch: langwatch,
	}
}

// GetMetadata passes through to the inner provider
func (p *langwatchProviderDecorator) GetMetadata() ProviderMetadata {
	return p.inner.GetMetadata()
}

// ListModels passes through to the inner provider
func (p *langwatchProviderDecorator) ListModels(ctx context.Context) ([]Model, error) {
	return p.inner.ListModels(ctx)
}

// NewChatStream creates a chat stream with langwatch tracking
func (p *langwatchProviderDecorator) NewChatStream(ctx context.Context, params ChatStreamParams) (ChatStreamIterator, error) {
	ctx, trace := p.langwatch.CreateTrace(ctx, params.ThreadID, langwatch.WithTraceMetadata(
		langwatch.Attribute(langwatch.AttributeNameUserID, params.ThreadOwnerID),
	))

	ctx, span := trace.StartLLMSpan(ctx, params.MessageID, langwatch.StartLLMSpanParams{
		Vendor: string(p.inner.GetMetadata().ProviderID),
		Model:  params.ModelID,
		Input: langwatch.StartLLMSpanChatMessageInput{
			Messages: convertToLangwatchMessage(params.Messages),
		},
	})

	iter, err := p.inner.NewChatStream(ctx, params)
	if err != nil {
		span.AddError(ctx, err)
		span.End(ctx)

		return nil, err
	}

	return &langwatchChatStreamIterator{
		inner:   iter,
		trace:   trace,
		llmSpan: span,
	}, nil
}

// langwatchChatStreamIterator wraps a ChatStreamIterator with langwatch capabilities
type langwatchChatStreamIterator struct {
	inner   ChatStreamIterator
	trace   langwatch.Trace
	llmSpan langwatch.SpanLLM
	content string
	ended   bool
}

func (i *langwatchChatStreamIterator) Next() bool {
	if i.ended {
		return false
	}

	hasNext := i.inner.Next()
	if !hasNext {
		// Stream is done, finalize the trace if there's no error
		if i.inner.Err() == nil {
			i.llmSpan.AddOutput(context.Background(), i.Content())
			i.llmSpan.End(context.Background())
		}

		i.ended = true
	}

	return hasNext
}

func (i *langwatchChatStreamIterator) Current() *ChatStreamEvent {
	return i.inner.Current()
}

func (i *langwatchChatStreamIterator) Content() string {
	i.content = i.inner.Content()

	return i.content
}

func (i *langwatchChatStreamIterator) Err() error {
	err := i.inner.Err()
	if err != nil && !i.ended {
		i.llmSpan.AddError(context.Background(), err)
		i.llmSpan.End(context.Background())

		i.ended = true
	}
	return err
}

func convertToLangwatchMessage(messages []Message) []langwatch.LLMSpanInputChatMessage {
	converted := make([]langwatch.LLMSpanInputChatMessage, len(messages))
	for i, m := range messages {
		converted[i] = langwatch.LLMSpanInputChatMessage{
			Role:    string(m.Role),
			Message: m.Content,
		}
	}

	return converted
}
