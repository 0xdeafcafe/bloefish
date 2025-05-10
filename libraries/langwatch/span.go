package langwatch

import (
	"encoding/json"
	"log"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type SpanType string

const (
	SpanTypeSpan       SpanType = "span"
	SpanTypeLLM        SpanType = "llm"
	SpanTypeChain      SpanType = "chain"
	SpanTypeTool       SpanType = "tool"
	SpanTypeAgent      SpanType = "agent"
	SpanTypeGuardrail  SpanType = "guardrail"
	SpanTypeEvaluation SpanType = "evaluation"
	SpanTypeRAG        SpanType = "rag"
	SpanTypeWorkflow   SpanType = "workflow"
	SpanTypeComponent  SpanType = "component"
	SpanTypeModule     SpanType = "module"
	SpanTypeServer     SpanType = "server"
	SpanTypeClient     SpanType = "client"
	SpanTypeProducer   SpanType = "producer"
	SpanTypeConsumer   SpanType = "consumer"
	SpanTypeTask       SpanType = "task"
	SpanTypeUnknown    SpanType = "unknown"
)

type SpanMetrics struct {
	PromptTokens     *int     `json:"prompt_tokens"`
	CompletionTokens *int     `json:"completion_tokens"`
	Cost             *float64 `json:"cost"`
}

type SpanTimestamps struct {
	StartedAtUnix    int64  `json:"started_at"`
	FirstTokenAtUnix *int64 `json:"first_token_at"`
	FinishedAtUnix   int64  `json:"finished_at"`
}

type SpanRAGContextChunk struct {
	DocumentID string `json:"document_id"`
	ChunkID    string `json:"chunk_id"`
	Content    any    `json:"content"`
}

type LangWatchSpan struct {
	trace.Span
}

func (s *LangWatchSpan) RecordInput(input any) {
	jsonStr, err := json.Marshal(input)
	if err != nil {
		log.Default().Printf("error marshalling input: %v", err)
	}

	s.SetAttributes(attribute.String(AttributeLangWatchInputKey, string(jsonStr)))
}

func (s *LangWatchSpan) RecordOutput(output any) {
	jsonStr, err := json.Marshal(output)
	if err != nil {
		log.Default().Printf("error marshalling output: %v", err)
	}

	s.SetAttributes(attribute.String(AttributeLangWatchOutputKey, string(jsonStr)))
}

func (s *LangWatchSpan) SetType(spanType SpanType) {
	s.SetAttributes(attribute.String(AttributeLangWatchSpanTypeKey, string(spanType)))
}

func (s *LangWatchSpan) SetRequestModel(model string) {
	s.SetAttributes(attribute.String(string(semconv.GenAiRequestModelKey), model))
}

func (s *LangWatchSpan) SetResponseModel(model string) {
	s.SetAttributes(attribute.String(string(semconv.GenAiResponseModelKey), model))
}

func (s *LangWatchSpan) SetMetrics(metrics SpanMetrics) {
	jsonStr, err := json.Marshal(metrics)
	if err != nil {
		log.Default().Printf("error marshalling metrics: %v", err)
	}

	s.SetAttributes(attribute.String(AttributeLangWatchMetricsKey, string(jsonStr)))
}

func (s *LangWatchSpan) SetTimestamps(timestamps SpanTimestamps) {
	jsonStr, err := json.Marshal(timestamps)
	if err != nil {
		log.Default().Printf("error marshalling timestamps: %v", err)
	}

	s.SetAttributes(attribute.String(AttributeLangWatchTimestampsKey, string(jsonStr)))
}

func (s *LangWatchSpan) SetRAGContextChunks(contexts []SpanRAGContextChunk) {
	jsonStr, err := json.Marshal(contexts)
	if err != nil {
		log.Default().Printf("error marshalling contexts: %v", err)
	}

	s.SetAttributes(attribute.String(AttributeLangWatchRAGContextsKey, string(jsonStr)))
}

func (s *LangWatchSpan) SetRAGContextChunk(context SpanRAGContextChunk) {
	s.SetRAGContextChunks([]SpanRAGContextChunk{context})
}
