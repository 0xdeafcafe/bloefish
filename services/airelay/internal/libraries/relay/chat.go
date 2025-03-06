package relay

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

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

type ChatStreamParams struct {
	ThreadID      string
	ThreadOwnerID string
	MessageID     string
	ModelID       string
	Messages      []Message
}

type ChatStreamEvent struct {
	Content string
	Done    bool
}

type ChatStreamIterator interface {
	Next() bool
	Current() *ChatStreamEvent
	Content() string
	Err() error
}
