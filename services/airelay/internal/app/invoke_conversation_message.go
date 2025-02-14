package app

import (
	"context"
	"fmt"

	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/openai/openai-go"
)

func (a *App) InvokeConversationMessage(ctx context.Context, req *airelay.InvokeConversationMessageRequest) (*airelay.InvokeConversationMessageResponse, error) {
	chatCompletionMessages, err := a.prepareOpenAIChatCompletionMessages(ctx, req.AIRelayOptions, req.Messages)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare OpenAI chat completion messages: %w", err)
	}

	completion, err := a.OpenAI.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(chatCompletionMessages),
		Seed:     openai.Int(1),
		Model:    openai.F(req.AIRelayOptions.ModelID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create completion: %w", err)
	}

	if len(completion.Choices) == 0 {
		return nil, fmt.Errorf("no choices in completion")
	}

	return &airelay.InvokeConversationMessageResponse{
		MessageContent: completion.Choices[0].Message.Content,
	}, nil
}
