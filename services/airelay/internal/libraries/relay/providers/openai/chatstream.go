package openai

import (
	"context"

	"github.com/0xdeafcafe/bloefish/services/airelay/internal/libraries/relay"

	"github.com/openai/openai-go"
	oaiClient "github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/ssestream"
)

type openAIChatStreamIterator struct {
	inner    *ssestream.Stream[oaiClient.ChatCompletionChunk]
	acc      oaiClient.ChatCompletionAccumulator
	current  *relay.ChatStreamEvent
	complete bool
}

func (i *openAIChatStreamIterator) Next() bool {
	if i.complete {
		return false
	}

	if !i.inner.Next() {
		i.complete = true
		return false
	}

	chunk := i.inner.Current()
	i.acc.AddChunk(chunk)

	i.current = &relay.ChatStreamEvent{
		Content: chunk.Choices[0].Delta.Content,
		Done:    false,
	}

	return true
}

func (i *openAIChatStreamIterator) Current() *relay.ChatStreamEvent {
	return i.current
}

func (i *openAIChatStreamIterator) Content() string {
	return i.acc.Choices[0].Message.Content
}

func (i *openAIChatStreamIterator) Err() error {
	return i.inner.Err()
}

func (p *Provider) NewChatStream(ctx context.Context, params relay.ChatStreamParams) (relay.ChatStreamIterator, error) {
	messages := make([]openai.ChatCompletionMessageParamUnion, len(params.Messages))
	for i, msg := range params.Messages {
		content := msg.Content
		switch msg.Role {
		case relay.RoleSystem:
			messages[i] = openai.AssistantMessage(content)
		case relay.RoleUser:
			messages[i] = openai.UserMessage(content)
		case relay.RoleAssistant:
			messages[i] = openai.AssistantMessage(content)
		}
	}

	stream := p.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(params.ModelID),
	})

	return &openAIChatStreamIterator{
		inner: stream,
	}, nil
}
