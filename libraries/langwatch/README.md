# langwatch-go

This is a Go client for the [Langwatch 🏰](https://langwatch.ai) API.

## Feature set

| Feature name | Supported |
| ------------ | --------- |
| LLM Spans    | Partial   |
| RAG Spans    | 🚫        |

## Usage

```go
package main

import (
	"github.com/0xdeafcafe/bloefish/libraries/langwatch"
	"github.com/0xdeafcafe/bloefish/libraries/ksuid"
	"github.com/0xdeafcafe/bloefish/libraries/openai"
)

var client langwatch.Client
var oaiClient openai.MockClient

func init() {
	ctx := context.Background()
	client = langwatch.NewClient()
	oaiClient = openai.NewMockClient()

	threadID := ksuid.Generate(ctx, "thread").String()
	userID := ksuid.Generate(ctx, "user").String()

	// Create a trace
	ctx, trace := client.CreateTrace(ctx, threadID, langwatch.WithTraceMetadata(
		langwatch.Attribute("user_id", userID),
		langwatch.Attribute("app_id", "AppID"),
		langwatch.Attribute("app_version", "1.0.0"),
	))
	defer trace.End()

	messages := createOpenAIMessages()

	ctx, span := trace.StartLLMSpan(ctx, langwatch.StartLLMSpanCommand{
		Name: "llm",
		Model: "gpt-4.0",
		Messages: langwatch.CoerceOpenAIMessages(messages),
		// Messages: langwatch.CoerceGrokMessages(messages),
		// Messages: langwatch.CoerceDeepseekMessages(messages),
	})
	defer span.End()

	span.AddOutput(ctx, langwatch.LLMSpanOutput{
		Output: 
	})

	// Create a span


}
```
