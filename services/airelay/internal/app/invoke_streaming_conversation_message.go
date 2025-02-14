package app

import (
	"context"
	"fmt"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/0xdeafcafe/bloefish/services/stream"
	"github.com/openai/openai-go"
)

func (a *App) InvokeStreamingConversationMessage(ctx context.Context, req *airelay.InvokeStreamingConversationMessageRequest) (*airelay.InvokeStreamingConversationMessageResponse, error) {
	chatCompletionMessages, err := a.prepareOpenAIChatCompletionMessages(ctx, req.AIRelayOptions, req.Messages)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare OpenAI chat completion messages: %w", err)
	}

	fmt.Println("a")

	chatStream := a.OpenAI.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(chatCompletionMessages),
		Seed:     openai.Int(1),
		Model:    openai.F(req.AIRelayOptions.ModelID),
	})
	acc := openai.ChatCompletionAccumulator{}

	// TODO(afr): Add an index to the stream message to help the client reconcile.
	// Currently it would be rare for the messages to be out of sync, as we wait for the
	// RPC call sending the message fragment to complete before we send the next message,
	// however it would be quicker if we send the RPC call in a go routine, allowing the
	// stream to continue processing and firing events to the client. The client would then
	// be responsible for reconciling the messages using the index.
	for chatStream.Next() {
		chunk := chatStream.Current()
		acc.AddChunk(chunk)

		fmt.Println("b")

		if content, ok := acc.JustFinishedContent(); ok {
			if err := a.StreamService.SendMessageFull(ctx, &stream.SendMessageFullRequest{
				ChannelID:      req.StreamingChannelID,
				MessageContent: content,
			}); err != nil {
				return nil, fmt.Errorf("failed to send message fragment: %w", err)
			}
		}

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			if err := a.StreamService.SendMessageFull(ctx, &stream.SendMessageFullRequest{
				ChannelID:      req.StreamingChannelID,
				MessageContent: refusal,
			}); err != nil {
				return nil, fmt.Errorf("failed to send message fragment: %w", err)
			}
		}

		if len(chunk.Choices) > 0 {
			clog.Get(ctx).WithField("chunk", chunk).Info("sending message fragment")

			if err := a.StreamService.SendMessageFragment(ctx, &stream.SendMessageFragmentRequest{
				ChannelID:      req.StreamingChannelID,
				MessageContent: chunk.Choices[0].Delta.Content,
			}); err != nil {
				return nil, fmt.Errorf("failed to send message fragment: %w", err)
			}
		}
	}

	if err := chatStream.Err(); err != nil {
		return nil, fmt.Errorf("failed to stream chat completions: %w", err)
	}

	return &airelay.InvokeStreamingConversationMessageResponse{
		MessageContent: acc.Choices[0].Message.Content,
	}, nil
}
