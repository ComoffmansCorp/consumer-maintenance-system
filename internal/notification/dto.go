package notification

import "time"

type DTO struct {
	ID        int64          `json:"id"`
	Type      string         `json:"type"`
	Title     string         `json:"title"`
	Message   string         `json:"message"`
	Payload   map[string]any `json:"payload,omitempty"`
	Read      bool           `json:"read"`
	CreatedAt time.Time      `json:"createdAt"`
}

func ToDTO(n Notification) DTO {
	return DTO{
		ID:        n.ID,
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		Payload:   n.Payload,
		Read:      n.ReadAt != nil,
		CreatedAt: n.CreatedAt,
	}
}
