package langwatch

import "context"

type Trace interface {
	StartLLMSpan(ctx context.Context, spanID string, params StartLLMSpanParams) (context.Context, SpanLLM)
}

type trace struct {
	attributes []attribute[any]
	client     *client
}

func (c *client) CreateTrace(ctx context.Context, theadID string, opts ...TraceOption) (context.Context, Trace) {
	t := &trace{
		attributes: []attribute[any]{},
		client:     c,
	}

	for _, opt := range opts {
		opt(t)
	}

	return ctx, t
}
