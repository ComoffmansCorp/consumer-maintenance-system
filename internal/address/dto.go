package address

import "time"

type DTO struct {
	ID           int64     `json:"id"`
	Street       string    `json:"street"`
	House        string    `json:"house"`
	Building     string    `json:"building,omitempty"`
	Apartment    string    `json:"apartment,omitempty"`
	ConsumerID   *int64    `json:"consumerId,omitempty"`
	ConsumerName string    `json:"consumerName,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type CreateRequest struct {
	Street     string `json:"street"`
	House      string `json:"house"`
	Building   string `json:"building"`
	Apartment  string `json:"apartment"`
	ConsumerID *int64 `json:"consumerId"`
}

type UpdateRequest struct {
	Street     string `json:"street"`
	House      string `json:"house"`
	Building   string `json:"building"`
	Apartment  string `json:"apartment"`
	ConsumerID *int64 `json:"consumerId"`
}

func ToDTO(a Address) DTO {
	return DTO{
		ID:         a.ID,
		Street:     a.Street,
		House:      a.House,
		Building:   a.Building,
		Apartment:  a.Apartment,
		ConsumerID: a.ConsumerID,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}
