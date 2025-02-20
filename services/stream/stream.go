package stream

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
)

type Service interface {
	SendMessageFull(context.Context, *SendMessageFullRequest) error
	SendMessageFragment(context.Context, *SendMessageFragmentRequest) error
	SendErrorMessage(context.Context, *SendErrorMessageRequest) error
}

type StreamedMessageType string

const (
	StreamedMessageTypeMessageFull     StreamedMessageType = "message_full"
	StreamedMessageTypeMessageFragment StreamedMessageType = "message_fragment"
	StreamedMessageTypeError           StreamedMessageType = "error"
)

type SendMessageFullRequest struct {
	ChannelID      string `json:"channel_id"`
	MessageContent string `json:"message_content"`
}

type SendMessageFragmentRequest struct {
	ChannelID      string `json:"channel_id"`
	MessageContent string `json:"message_content"`
}

type SendErrorMessageRequest struct {
	ChannelID string `json:"channel_id"`
	Error     cher.E `json:"error"`
}

type StreamedMessage struct {
	ChannelID string              `json:"channel_id"`
	Type      StreamedMessageType `json:"type"`

	MessageFull     *string `json:"message_full"`
	MessageFragment *string `json:"message_fragment"`
	Error           *cher.E `json:"error"`
}
