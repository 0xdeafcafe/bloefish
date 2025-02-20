package ollama

type Role string

const (
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
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
