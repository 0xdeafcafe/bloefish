package stream

import "context"

type Service interface {
	SendMessageFull(context.Context, *SendMessageFullRequest) error
	SendMessageFragment(context.Context, *SendMessageFragmentRequest) error

	SendStreamedMessage(context.Context, *StreamedMessage) error
}

type StreamedMessageType string

const (
	StreamedMessageTypeMessageFull     StreamedMessageType = "message_full"
	StreamedMessageTypeMessageFragment StreamedMessageType = "message_fragment"
)

type SendMessageFullRequest struct {
	ChannelID      string `json:"channel_id"`
	MessageContent string `json:"message_content"`
}

type SendMessageFragmentRequest struct {
	ChannelID      string `json:"channel_id"`
	MessageContent string `json:"message_content"`
}

type StreamedMessage struct {
	ChannelID       string              `json:"channel_id"`
	Type            StreamedMessageType `json:"type"`
	MessageFull     *string             `json:"message_full"`
	MessageFragment *string             `json:"message_fragment"`
}
