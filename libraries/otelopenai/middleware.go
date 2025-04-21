package otelopenai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path"

	oaioption "github.com/openai/openai-go/option"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv/v1.13.0/httpconv"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "github.com/0xdeafcafe/bloefish/libraries/otelopenai"
)

type traceware struct {
	config
	name    string
	tracer  oteltrace.Tracer
	handler http.Handler
}

// Middleware sets up a handler to start tracing the requests made to OpenAI by the
// OpenAI library. The name parameter should describe the purpose of this specific OpenAI
// client instance.
func Middleware(name string, opts ...Option) oaioption.Middleware {
	cfg := config{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	if cfg.tracerProvider == nil {
		cfg.tracerProvider = otel.GetTracerProvider()
	}
	tracer := cfg.tracerProvider.Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion("0.0.1"),
	)
	if cfg.propagators == nil {
		cfg.propagators = otel.GetTextMapPropagator()
	}

	return func(req *http.Request, next oaioption.MiddlewareNext) (*http.Response, error) {
		// Extract operation name from the URL path
		operation := path.Base(req.URL.Path)

		// Start a new span
		ctx := req.Context()
		spanName := "openai." + operation
		ctx, span := tracer.Start(ctx, spanName,
			oteltrace.WithAttributes(
				semconv.HTTPMethodKey.String(req.Method),
				semconv.HTTPURLKey.String(req.URL.String()),
				attribute.String("openai.operation", operation),
			),
			oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		)
		defer span.End()

		// Read and record request body
		var reqBody []byte
		if req.Body != nil {
			reqBody, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			// Record request content
			span.SetAttributes(attribute.String("openai.request", string(reqBody)))
		}

		// Call the next middleware
		resp, err := next(req.WithContext(ctx))
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return resp, err
		}

		// Read and record response body
		if resp != nil && resp.Body != nil {
			respBody, _ := io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

			// Record response content
			span.SetAttributes(attribute.String("openai.response", string(respBody)))

			// Try to extract and record token counts
			var respData map[string]interface{}
			if err := json.Unmarshal(respBody, &respData); err == nil {
				if usage, ok := respData["usage"].(map[string]interface{}); ok {
					if promptTokens, ok := usage["prompt_tokens"].(float64); ok {
						span.SetAttributes(attribute.Int("openai.tokens.prompt", int(promptTokens)))
					}
					if completionTokens, ok := usage["completion_tokens"].(float64); ok {
						span.SetAttributes(attribute.Int("openai.tokens.completion", int(completionTokens)))
					}
					if totalTokens, ok := usage["total_tokens"].(float64); ok {
						span.SetAttributes(attribute.Int("openai.tokens.total", int(totalTokens)))
					}
				}
			}
		}

		// Record response status
		if resp != nil {
			span.SetAttributes(semconv.HTTPStatusCode(resp.StatusCode))
			span.SetStatus(httpconv.ClientStatus(resp.StatusCode))
		}

		return resp, nil
	}
}
