package langwatch

import (
	"context"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/ptr"
)

type SpanLLM interface {
	AddOutput(ctx context.Context, output string)
	AddError(ctx context.Context, err error)
	End(ctx context.Context)
}

type LLMSpanInputChatMessage struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

type LLMSpanOutputChatMessage struct {
	Role         string        `json:"role"`
	Message      string        `json:"message"`
	FunctionCall interface{}   `json:"function_call"`
	ToolCalls    []interface{} `json:"tool_calls"`
}

type StartLLMSpanParams struct {
	Vendor string
	Model  string
	Name   *string
	Input  StartLLMSpanChatMessageInput
}

type StartLLMSpanChatMessageInput struct {
	Messages []LLMSpanInputChatMessage
}

type spanLLM struct {
	client     *client
	trace      Trace
	spanID     string
	name       *string
	vendor     *string
	model      *string
	input      *spanLLMInput
	output     *spanLLMOutput
	error      error
	timestamps *spanLLMTimestamps
	metrics    *spanLLMMetrics
}

type spanLLMInput struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
type spanLLMOutput struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
type spanLLMParams struct {
	Temperature *float64 `json:"temperature,omitempty"`
	Stream      *bool    `json:"stream,omitempty"`
}
type spanLLMTimestamps struct {
	StartedAt  *int64 `json:"started_at,omitempty"`
	FinishedAt *int64 `json:"finished_at,omitempty"`
}
type spanLLMMetrics struct {
	PromptTokens     *int `json:"prompt_tokens,omitempty"`
	CompletionTokens *int `json:"completion_tokens,omitempty"`
}

type spanCollectorRequest struct {
	TraceID  string                      `json:"trace_id"`
	Spans    []*spanCollectorRequestSpan `json:"spans"`
	Metadata map[string]string           `json:"metadata"`
}
type spanCollectorRequestSpan struct {
	Type       string
	SpanID     string
	Vendor     string
	Model      string
	Input      interface{}
	Output     interface{}
	Params     interface{}
	Metrics    interface{}
	Timestamps interface{}
}
type spanCollectorRequestMetadata struct {
}

func (t *trace) StartLLMSpan(ctx context.Context, spanID string, params StartLLMSpanParams) (context.Context, SpanLLM) {
	return ctx, &spanLLM{
		client: t.client,
		trace:  t,

		spanID: spanID,
		vendor: &params.Vendor,
		model:  &params.Model,

		input: &spanLLMInput{
			Type:  "chat_messages",
			Value: params.Input.Messages,
		},

		timestamps: &spanLLMTimestamps{
			StartedAt: ptr.P(time.Now().UnixNano()),
		},
	}
}

func (s *spanLLM) AddOutput(ctx context.Context, output string) {
	s.output = &spanLLMOutput{
		Type: "chat_messages",
		Value: []LLMSpanOutputChatMessage{{
			Role:         "bot",
			Message:      output,
			FunctionCall: nil,
			ToolCalls:    []interface{}{},
		}},
	}
	s.timestamps.FinishedAt = ptr.P(time.Now().UnixNano())
}

func (s *spanLLM) AddError(ctx context.Context, err error) {
	s.error = err
}

func (s *spanLLM) End(ctx context.Context) {
	s.client.Collect(ctx, nil)
}
