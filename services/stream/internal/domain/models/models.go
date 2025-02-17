package models

import "github.com/0xdeafcafe/bloefish/libraries/cher"

type StreamMessageType string

const (
	StreamMessageTypeMessageFull     StreamMessageType = "message_full"
	StreamMessageTypeMessageFragment StreamMessageType = "message_fragment"
	StreamMessageTypeError           StreamMessageType = "error"
)

type StreamMessage struct {
	ChannelID       string            `json:"channel_id"`
	Type            StreamMessageType `json:"type"`
	MessageFull     *string           `json:"message_full"`
	MessageFragment *string           `json:"message_fragment"`
	Error           *cher.E           `json:"error"`
}
