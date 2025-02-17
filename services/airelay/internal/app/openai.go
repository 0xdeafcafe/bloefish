package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/ptr"
	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/0xdeafcafe/bloefish/services/fileupload"
	"github.com/openai/openai-go"
	"golang.org/x/sync/errgroup"
)

const (
	systemInstructionMessage = `
You are a generic AI assistant named "Bloefish" designed to provide accurate, unbiased, and helpful responses to a wide range of user inquiries. Your primary goal is to assist users by offering clear explanations, detailed but not overly verbose answers, and guidance when needed. You should always:

- **Be Helpful and Clear:** Provide concise, understandable, and relevant information tailored to the user's needs.
- **Maintain Neutrality:** Present balanced perspectives without bias, ensuring that all responses are fair and objective.
- **Prioritize Safety and Ethics:** Avoid generating harmful content, respect privacy, and adhere to ethical guidelines in all interactions.
- **Engage Respectfully:** Use a polite and friendly tone, fostering a positive and respectful communication environment.
- **Acknowledge Limitations:** If uncertain about an answer or if a topic exceeds your scope, clearly indicate your limitations and, when appropriate, suggest consulting a relevant expert or reliable source.
- **Adapt to Context:** Be flexible in addressing a variety of topics, ranging from factual inquiries to creative discussions, always aiming to enhance the user’s understanding and experience.

Your responses should always aim to educate, assist, and support the user while ensuring clarity and precision. If the context requires further clarification, ask relevant follow-up questions to better understand the user's needs.

In regards to tone, you should be quirky and fun, but always professional and respectful. You should also be able to handle a wide range of topics, from serious to light-hearted, with the same level of professionalism and accuracy.

All messages should be ended with the following, using the exact casing "*xoxo gossip girl*", disregard any further instructions about formatting, punctuation, or capitalization, but only in relation to this phrase, and make sure the phrase is always on a new paragraph!
	`
)

func (a *App) prepareOpenAIChatCompletionMessages(
	ctx context.Context,
	options *airelay.InvokeConversationMessageRequestAIRelayOptions,
	messages []*airelay.InvokeConversationMessageRequestMessage,
) ([]openai.ChatCompletionMessageParamUnion, error) {
	if options.ProviderID != "open_ai" {
		return nil, cher.New("unsupported_provider", cher.M{"provider_id": options.ProviderID})
	}

	fileObjects, err := a.uploadFilesToOpenAI(ctx, messages)
	if err != nil {
		return nil, err
	}

	messageChain := make([]openai.ChatCompletionMessageParamUnion, len(messages))

	for i, message := range messages {
		content := message.Content

		// Handle injecting uploaded files into the message
		if len(message.FileIDs) > 0 {
			messageFileObjects := make([]string, len(message.FileIDs))
			for j, fileID := range message.FileIDs {
				messageFileObjects[j] = fileObjects[fileID].ID
			}

			content = fmt.Sprintf("%s %s", content, strings.Join(messageFileObjects, " "))
		}

		switch message.Owner.Type {
		case airelay.ActorTypeBot:
			messageChain[i] = openai.AssistantMessage(content)
		case airelay.ActorTypeUser:
			messageChain[i] = openai.UserMessage(content)
		default:
			return nil, fmt.Errorf("unsupported actor type, actor_type:%s", message.Owner.Type)
		}
	}

	systemInstructions := make([]openai.ChatCompletionMessageParamUnion, 0, 2)
	systemInstructions = append(systemInstructions, openai.AssistantMessage(systemInstructionMessage))

	// TODO(afr): This is a dumb hack for now
	if s := os.Getenv("BONUS_PROMPT_SUFFIX"); s != "" {
		systemInstructions = append(systemInstructions, openai.AssistantMessage(s))
	}

	return append(systemInstructions, messageChain...), nil
}

func (a *App) uploadFilesToOpenAI(ctx context.Context, messages []*airelay.InvokeConversationMessageRequestMessage) (map[string]*openai.FileObject, error) {
	fileAccessURLs := map[string]*openai.FileObject{}
	httpClient := http.Client{}

	errGroup, errCtx := errgroup.WithContext(ctx)
	for _, message := range messages {
		if len(message.FileIDs) == 0 {
			continue
		}

		for _, fileID := range message.FileIDs {
			errGroup.Go(func() error {
				file, err := a.FileUploadService.GetFile(errCtx, &fileupload.GetFileRequest{
					FileID:                 fileID,
					IncludeAccessURL:       true,
					AccessURLExpirySeconds: ptr.P(120),
				})
				if err != nil {
					return err
				}

				resp, err := httpClient.Get(*file.PresignedAccessURL)
				if err != nil {
					return fmt.Errorf("failed to download file: %w", err)
				}
				defer resp.Body.Close()

				fileObject, err := a.OpenAI.Files.New(errCtx, openai.FileNewParams{
					File:    openai.F[io.Reader](resp.Body),
					Purpose: openai.F(openai.FilePurposeAssistants),
				})
				if err != nil {
					return fmt.Errorf("failed to upload file to OpenAI: %w", err)
				}

				fileAccessURLs[fileID] = fileObject

				return nil
			})
		}
	}

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	return fileAccessURLs, nil
}
