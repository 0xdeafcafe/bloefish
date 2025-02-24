package relay

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/ssestream"
)

type openAIChatStreamIterator struct {
	inner    *ssestream.Stream[openai.ChatCompletionChunk]
	acc      openai.ChatCompletionAccumulator
	current  *ChatStreamEvent
	complete bool
}

func (c *Client) newOpenAIChatStream(ctx context.Context, params ChatStreamParams) (ChatStreamIterator, error) {
	messages := make([]openai.ChatCompletionMessageParamUnion, len(params.Messages))
	for i, msg := range params.Messages {
		content := msg.Content
		switch msg.Role {
		case RoleSystem:
			messages[i] = openai.AssistantMessage(content)
		case RoleUser:
			messages[i] = openai.UserMessage(content)
		case RoleAssistant:
			messages[i] = openai.AssistantMessage(content)
		}
	}

	stream := c.openAIClient.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(params.ModelID),
	})

	return &openAIChatStreamIterator{
		inner: stream,
	}, nil
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

	i.current = &ChatStreamEvent{
		Content: chunk.Choices[0].Delta.Content,
		Done:    false,
	}

	return true
}

func (i *openAIChatStreamIterator) Current() *ChatStreamEvent {
	return i.current
}

func (i *openAIChatStreamIterator) Content() string {
	return i.acc.Choices[0].Message.Content
}

func (i *openAIChatStreamIterator) Err() error {
	return i.inner.Err()
}
