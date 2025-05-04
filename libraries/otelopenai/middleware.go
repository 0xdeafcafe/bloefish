package otelopenai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path"
	"strings"

	oaioption "github.com/openai/openai-go/option"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerName             = "github.com/0xdeafcafe/bloefish/libraries/otelopenai"
	instrumentationVersion = "0.0.1" // TODO: Consider linking this to package version
)

// Middleware sets up a handler to start tracing the requests made to OpenAI by the
// OpenAI library.
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
		oteltrace.WithInstrumentationVersion(instrumentationVersion),
		oteltrace.WithSchemaURL(semconv.SchemaURL),
	)
	if cfg.propagators == nil {
		cfg.propagators = otel.GetTextMapPropagator()
	}

	return func(req *http.Request, next oaioption.MiddlewareNext) (*http.Response, error) {
		operation := path.Base(req.URL.Path)

		ctx := req.Context()
		spanName := "openai." + operation
		ctx, span := tracer.Start(ctx, spanName,
			oteltrace.WithAttributes(
				semconv.HTTPRequestMethodKey.String(req.Method),
				semconv.ServerAddressKey.String(req.URL.Hostname()),
				semconv.URLPathKey.String(req.URL.Path),
				attributeGenAISystem.String(openAISystemValue), // Defined in attributes.go
				attributeGenAIOperation.String(operation),
			),
			oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		)
		defer span.End()

		var reqBody []byte
		var isStreaming bool
		if req.Body != nil && req.Body != http.NoBody {
			var errRead error
			reqBody, errRead = io.ReadAll(req.Body)
			// Restore the body so the downstream handler can read it
			req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			if errRead == nil {
				if cfg.recordInput {
					// TODO(afr): Adjust this based on the operation
					span.SetAttributes(attributeLangwatchInputValue.String(string(reqBody)))
				}
				var reqData jsonData
				if err := json.Unmarshal(reqBody, &reqData); err == nil {
					if model, ok := getString(reqData, "model"); ok {
						span.SetAttributes(attributeGenAIRequestModel.String(model))
					}
					if temp, ok := getFloat64(reqData, "temperature"); ok {
						span.SetAttributes(attributeGenAIRequestTemperature.Float64(temp))
					}
					if topP, ok := getFloat64(reqData, "top_p"); ok {
						span.SetAttributes(attributeGenAIRequestTopP.Float64(topP))
					}
					if freqPenalty, ok := getFloat64(reqData, "frequency_penalty"); ok {
						span.SetAttributes(attributeGenAIRequestFrequencyPenalty.Float64(freqPenalty))
					}
					if presPenalty, ok := getFloat64(reqData, "presence_penalty"); ok {
						span.SetAttributes(attributeGenAIRequestPresencePenalty.Float64(presPenalty))
					}
					if maxTokens, ok := getInt(reqData, "max_tokens"); ok {
						span.SetAttributes(attributeGenAIRequestMaxTokens.Int(maxTokens))
					}

					if streamVal, ok := reqData["stream"].(bool); ok && streamVal {
						isStreaming = true
						span.SetAttributes(attributeGenAIRequestStream.Bool(true))
					} else {
						span.SetAttributes(attributeGenAIRequestStream.Bool(false))
					}
				} else {
					span.AddEvent("Failed to parse OpenAI request body JSON", oteltrace.WithAttributes(attribute.String("error", err.Error())))
				}
			} else {
				span.SetStatus(codes.Error, "failed to read request body")
				span.RecordError(errRead)
			}
		}

		resp, err := next(req.WithContext(ctx))
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			if resp != nil {
				span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(resp.StatusCode))
			}
			return resp, err
		}

		if resp != nil {
			span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(resp.StatusCode))
			if resp.StatusCode >= 400 {
				span.SetStatus(codes.Error, http.StatusText(resp.StatusCode))
				// Potentially read error body here if needed, even for non-JSON
			} else {
				span.SetStatus(codes.Ok, "")
			}

			// Streaming response bodies are not processed by this middleware.
			// The `openai-go` client library's stream consumption model prevents
			// reliable interception and parsing of stream chunks at the HTTP
			// middleware layer. Therefore, response attributes derived from stream
			// content (e.g., token usage, finish reasons, recorded output) are not
			// captured for streaming calls.
			// Only non-streaming JSON responses are fully parsed.
			if !isStreaming && resp.Body != nil && resp.Body != http.NoBody {
				// Handle non-streaming response body
				respBody, errRead := io.ReadAll(resp.Body)
				// Restore the *response* body so the client can read it
				resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
				if errRead == nil {
					contentType := resp.Header.Get("Content-Type")
					if strings.HasPrefix(contentType, "application/json") {
						if cfg.recordOutput {
							span.SetAttributes(attributeLangwatchOutputValue.String(string(respBody)))
						}
						var respData jsonData
						if err := json.Unmarshal(respBody, &respData); err == nil {
							setNonStreamResponseAttributes(span, respData)
						} else {
							span.AddEvent("Failed to parse non-stream OpenAI response body JSON", oteltrace.WithAttributes(attribute.String("error", err.Error())))
						}
					}
				} else {
					span.AddEvent("Failed to read non-stream OpenAI response body", oteltrace.WithAttributes(attribute.String("error", errRead.Error())))
				}
			}
		}

		return resp, nil
	}
}

// setNonStreamResponseAttributes extracts attributes from a standard JSON response body.
func setNonStreamResponseAttributes(span oteltrace.Span, respData jsonData) {
	if id, ok := getString(respData, "id"); ok {
		span.SetAttributes(attributeGenAIResponseID.String(id))
	}
	if model, ok := getString(respData, "model"); ok {
		span.SetAttributes(attributeGenAIResponseModel.String(model))
	}
	if sysFingerprint, ok := getString(respData, "system_fingerprint"); ok {
		span.SetAttributes(attributeGenAIOpenAIResponseSysFinger.String(sysFingerprint))
	}
	if usage, ok := respData["usage"].(jsonData); ok {
		if promptTokens, ok := getInt(usage, "prompt_tokens"); ok {
			span.SetAttributes(attributeGenAIUsageInputTokens.Int(promptTokens))
		}
		if completionTokens, ok := getInt(usage, "completion_tokens"); ok {
			span.SetAttributes(attributeGenAIUsageOutputTokens.Int(completionTokens))
		}
	}
	if choices, ok := respData["choices"].([]interface{}); ok {
		finishReasons := make([]string, 0, len(choices))
		for _, choiceRaw := range choices {
			if choice, ok := choiceRaw.(jsonData); ok {
				if reason, ok := getString(choice, "finish_reason"); ok {
					finishReasons = append(finishReasons, reason)
				}
			}
		}
		if len(finishReasons) > 0 {
			span.SetAttributes(attributeGenAIResponseFinishReasons.StringSlice(finishReasons))
		}
	}
}
