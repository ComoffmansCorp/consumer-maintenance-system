package consumer

import "time"

type DTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Type        Type      `json:"type"`
	Description string    `json:"description,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateRequest struct {
	Name        string `json:"name"`
	Type        Type   `json:"type"`
	Description string `json:"description"`
}

type UpdateRequest struct {
	Name        string `json:"name"`
	Type        Type   `json:"type"`
	Description string `json:"description"`
}

func ToDTO(c Consumer) DTO {
	return DTO{
		ID:          c.ID,
		Name:        c.Name,
		Type:        c.Type,
		Description: c.Description,
		Active:      c.Active,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
