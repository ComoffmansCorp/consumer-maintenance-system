package review

import "time"

type ReviewDTO struct {
	ID        int64     `json:"id"`
	RequestID int64     `json:"requestId"`
	ClientID  int64     `json:"clientId"`
	MasterID  int64     `json:"masterId"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateReviewRequest struct {
	RequestID int64  `json:"requestId"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}

func ToReviewDTO(r Review) ReviewDTO {
	return ReviewDTO{
		ID:        r.ID,
		RequestID: r.RequestID,
		ClientID:  r.ClientID,
		MasterID:  r.MasterID,
		Rating:    r.Rating,
		Comment:   r.Comment,
		CreatedAt: r.CreatedAt,
	}
}
