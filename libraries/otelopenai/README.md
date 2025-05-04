# otelopenai

This package provides OpenTelemetry instrumentation middleware for the official `openai-go` client library (`github.com/openai/openai-go`).

It automatically creates client spans for OpenAI API calls made through the instrumented client, adding relevant request and response attributes according to OpenTelemetry GenAI semantic conventions.

## Installation

```bash
go get github.com/0xdeafcafe/bloefish/libraries/otelopenai
```

## Usage

Wrap the `otelopenai.Middleware` using `option.WithMiddleware` when creating the `openai.Client`:

```go
package main

import (
	"context"
	"log"

	"github.com/0xdeafcafe/bloefish/libraries/otelopenai"
	"github.com/openai/openai-go"
	oaioption "github.com/openai/openai-go/option"
	"go.opentelemetry.io/otel"
)

func main() {
	// Configure OpenTelemetry provider (e.g., using Jaeger, OTLP exporter)
	// ... setup code for your tracer provider ...
	tracerProvider := otel.GetTracerProvider() // Get configured global provider

	// Create instrumented OpenAI client
	client := openai.NewClient(
		oaioption.WithAPIKey("YOUR_API_KEY"),
		oaioption.WithMiddleware(otelopenai.Middleware("my-openai-client",
			// Optional: Provide specific tracer provider
			otelopenai.WithTracerProvider(tracerProvider),
			// Optional: Capture request/response bodies (be mindful of sensitive data)
			otelopenai.WithCaptureInput(),
			otelopenai.WithCaptureOutput(),
		)),
	)

	// Make API calls as usual
	_, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.F(openai.ChatModelGPT4oMini),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Hello, OpenAI!"),
		}),
	})
	if err != nil {
		log.Fatalf("Chat completion failed: %v", err)
	}

	// ... shutdown tracer provider ...
}

```

## Configuration Options

The `Middleware` function accepts optional configuration functions:

*   `WithTracerProvider(provider oteltrace.TracerProvider)`: Specifies the OTel `TracerProvider`. Defaults to the global provider.
*   `WithPropagators(propagators propagation.TextMapPropagator)`: Specifies OTel propagators. Defaults to global propagators.
*   `WithTraceIDResponseHeader(key string)`: Adds the trace ID to the HTTP response header with the given `key`.
*   `WithTraceSampledResponseHeader(key string)`: Adds the trace sampling decision to the HTTP response header with the given `key`.
*   `WithCaptureInput()`: Records the full HTTP request body as the `langwatch.input.value` span attribute. Use with caution if requests contain sensitive data.
*   `WithCaptureOutput()`: Records the full HTTP response body (for non-streaming responses) as the `langwatch.output.value` span attribute. Use with caution if responses contain sensitive data. **Note:** See Streaming Limitation below.

## Collected Attributes

The middleware adds attributes to the client span, following [OpenTelemetry GenAI Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/gen-ai/general/) where applicable.

**Request Attributes:**

*   `gen_ai.system` (=`openai`)
*   `gen_ai.request.model`
*   `gen_ai.request.temperature`
*   `gen_ai.request.top_p`
*   `gen_ai.request.frequency_penalty`
*   `gen_ai.request.presence_penalty`
*   `gen_ai.request.max_tokens`
*   `gen_ai.request.stream` (boolean)
*   `gen_ai.operation.name` (e.g., `completions`)
*   `langwatch.input.value` (if `WithCaptureInput()` is used)

**Response Attributes (Non-Streaming Only):**

*   `gen_ai.response.id`
*   `gen_ai.response.model`
*   `gen_ai.response.finish_reasons`
*   `gen_ai.usage.input_tokens`
*   `gen_ai.usage.output_tokens`
*   `gen_ai.openai.response.system_fingerprint`
*   `langwatch.output.value` (if `WithCaptureOutput()` is used)

Standard HTTP client attributes (`http.request.method`, `url.path`, `server.address`, `http.response.status_code`) are also included.

## Streaming Limitation

Due to the internal implementation of stream handling in the `openai-go` client library, this HTTP-level middleware **cannot reliably parse the content of streaming responses**. 

Therefore, for streaming API calls:

*   Request attributes (including `gen_ai.request.stream=true`) are captured correctly.
*   Response attributes derived from the stream's *content* (e.g., `gen_ai.response.id`, `gen_ai.response.model`, `gen_ai.usage.output_tokens`, `gen_ai.response.finish_reasons`) are **not** captured by this middleware.
*   The `WithCaptureOutput()` option **does not** record the accumulated stream content.

Only non-streaming JSON responses are fully parsed for response attributes.
