package chat

import "time"

type Message struct {
	ID        int64
	RequestID int64
	SenderID  int64
	Text      string
	CreatedAt time.Time
	ReadAt    *time.Time
}
