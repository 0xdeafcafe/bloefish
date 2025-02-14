package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/ptr"
	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/0xdeafcafe/bloefish/services/fileupload"
	"github.com/openai/openai-go"
	"golang.org/x/sync/errgroup"
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

	chatCompletionMessages := make([]openai.ChatCompletionMessageParamUnion, len(messages))
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
			chatCompletionMessages[i] = openai.AssistantMessage(content)
		case airelay.ActorTypeUser:
			chatCompletionMessages[i] = openai.UserMessage(content)
		default:
			return nil, fmt.Errorf("unsupported actor type, actor_type:%s", message.Owner.Type)
		}
	}

	return chatCompletionMessages, nil
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
