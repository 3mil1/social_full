package dto

type ChatMessage struct {
	From    	string    `json:"from,"`
	Content 	string    `json:"content,"`
	CreatedAt 	string 	`json:"created_at,omitempty"`
}
