package relay

import "errors"

var ErrUnsupportedProvider = errors.New("unsupported provider")

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Provider struct {
	ID     string
	Name   string
	Models []Model
}

type Model struct {
	ID          string
	Name        string
	Description string
}

type Message struct {
	Role    Role
	Content string
}

func NewChatAssistantMessage(content string) Message {
	return Message{
		Role:    RoleAssistant,
		Content: content,
	}
}

func NewChatSystemMessage(content string) Message {
	return Message{
		Role:    RoleSystem,
		Content: content,
	}
}

func NewChatUserMessage(content string) Message {
	return Message{
		Role:    RoleUser,
		Content: content,
	}
}
