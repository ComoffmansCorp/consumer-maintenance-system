package chat

import "time"

type MessageDTO struct {
	ID        int64      `json:"id"`
	RequestID int64      `json:"requestId"`
	SenderID  int64      `json:"senderId"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"createdAt"`
	ReadAt    *time.Time `json:"readAt,omitempty"`
}

type SendMessageRequest struct {
	Text string `json:"text"`
}

func ToMessageDTO(m Message) MessageDTO {
	return MessageDTO(m)
}
