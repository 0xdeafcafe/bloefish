package ollama

type StreamingChatEvent struct {
	Model        string `json:"model"`
	CreatedAtStr string `json:"created_at"` // This is purposefully a string to avoid marshalling time`
	Done         bool   `json:"done"`

	// The following are only set when `done` = false
	Message *StreamingChatEventMessage `json:"message"`

	// The following are only set when `done` = true
	TotalDuration      int `json:"total_duration"`
	LoadDuration       int `json:"load_duration"`
	PromptEvalCount    int `json:"prompt_eval_count"`
	PromptEvalDuration int `json:"prompt_eval_duration"`
	EvalCount          int `json:"eval_count"`
	EvalDuration       int `json:"eval_duration"`
}

type StreamingChatEventMessage struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}
