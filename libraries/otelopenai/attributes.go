package otelopenai

import "go.opentelemetry.io/otel/attribute"

const (
	// Otel Spec GenAI Attributes
	attributeGenAISystem                  = attribute.Key("gen_ai.system")
	attributeGenAIRequestModel            = attribute.Key("gen_ai.request.model")
	attributeGenAIRequestTemperature      = attribute.Key("gen_ai.request.temperature")
	attributeGenAIRequestTopP             = attribute.Key("gen_ai.request.top_p")
	attributeGenAIRequestFrequencyPenalty = attribute.Key("gen_ai.request.frequency_penalty")
	attributeGenAIRequestPresencePenalty  = attribute.Key("gen_ai.request.presence_penalty")
	attributeGenAIRequestMaxTokens        = attribute.Key("gen_ai.request.max_tokens")
	attributeGenAIRequestStream           = attribute.Key("gen_ai.request.stream")
	attributeGenAIUsageInputTokens        = attribute.Key("gen_ai.usage.input_tokens")
	attributeGenAIUsageOutputTokens       = attribute.Key("gen_ai.usage.output_tokens")
	attributeGenAIResponseID              = attribute.Key("gen_ai.response.id")
	attributeGenAIResponseModel           = attribute.Key("gen_ai.response.model")
	attributeGenAIResponseFinishReasons   = attribute.Key("gen_ai.response.finish_reasons") // String Slice
	attributeGenAIOperation               = attribute.Key("gen_ai.operation.name")
	attributeGenAIOpenAIResponseSysFinger = attribute.Key("gen_ai.openai.response.system_fingerprint")

	// LangWatch specific attributes
	attributeLangwatchInputValue  = attribute.Key("langwatch.input.value")
	attributeLangwatchOutputValue = attribute.Key("langwatch.output.value")
)

const (
	openAISystemValue = "openai"
)
