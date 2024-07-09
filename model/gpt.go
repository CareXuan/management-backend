package model

type GptResponse struct {
	Id      string      `json:"id"`
	Object  string      `json:"object"`
	Created int         `json:"created"`
	Model   string      `json:"model"`
	Choices []GptChoice `json:"choices"`
	Usage   GptUsage    `json:"usage"`
}

type GptChoice struct {
	Index        int              `json:"index"`
	Message      GptChoiceMessage `json:"message"`
	FinishReason string           `json:"finish_reason"`
}

type GptChoiceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GptUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type GptQuestion struct {
	Id          int    `json:"id" xorm:"pk autoincr INT(11)"`
	UserId      int    `json:"user_id" xorm:"INT(11) not null default 0"`
	Question    string `json:"question" xorm:"text not null"`
	Answer      string `json:"answer" xorm:"text not null"`
	CreatedAt   int    `json:"created_at" xorm:"INT(11) not null default 0"`
	CreatedTime string `json:"created_time" xorm:"-"`
}
