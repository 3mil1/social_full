package entity

type ChatMessage struct {
	ChatID    int    ` json:"chat_id,omitempty"`
	MessageID int    `json:"message_id,omitempty"`
	From      string `json:"from,omitempty"`
	Content   string `json:"content,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}


type GroupChatMessage struct {
	GroupId		int		`json:"group_id,omitempty"`
	MessageID int    `json:"message_id,omitempty"`
	FromId      string `json:"from_id,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	LatsName      string `json:"last_name,omitempty"`
	Content   string `json:"content,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}