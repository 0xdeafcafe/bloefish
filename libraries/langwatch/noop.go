package langwatch

import "context"

type noopClient struct{}
type noopTrace struct{}
type noopSpanLLM struct{}

func NewNoopClient() Client {
	return &noopClient{}
}

func (c *noopClient) CreateTrace(ctx context.Context, theadID string, _ ...TraceOption) (context.Context, Trace) {
	t := &noopTrace{}

	return ctx, t
}

func (t *noopTrace) StartLLMSpan(ctx context.Context, spanID string, params StartLLMSpanParams) (context.Context, SpanLLM) {
	return ctx, &noopSpanLLM{}
}

func (s *noopSpanLLM) AddOutput(ctx context.Context, output string) {}
func (s *noopSpanLLM) AddError(ctx context.Context, err error)      {}
func (s *noopSpanLLM) End(ctx context.Context)                      {}
