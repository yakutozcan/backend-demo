package models

type Message struct {
	ID        int    `json:"id"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	Sent      bool   `json:"sent"`
	CreatedAt string `json:"created_at"`
}
