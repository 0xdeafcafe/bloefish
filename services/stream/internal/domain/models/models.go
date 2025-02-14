package models

type StreamMessageType string

const (
	StreamMessageTypeMessageFull     StreamMessageType = "message_full"
	StreamMessageTypeMessageFragment StreamMessageType = "message_fragment"
)

type StreamMessage struct {
	ChannelID       string            `json:"channel_id"`
	Type            StreamMessageType `json:"type"`
	MessageFull     *string           `json:"message_full"`
	MessageFragment *string           `json:"message_fragment"`
}
