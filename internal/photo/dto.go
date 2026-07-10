package photo

import "time"

type DTO struct {
	ID               int64     `json:"id"`
	Note             string    `json:"note,omitempty"`
	OriginalFilename string    `json:"originalFilename"`
	ContentType      string    `json:"contentType"`
	SizeBytes        int64     `json:"sizeBytes"`
	CreatedAt        time.Time `json:"createdAt"`
}

func ToDTO(p Photo) DTO {
	return DTO{
		ID:               p.ID,
		Note:             p.Note,
		OriginalFilename: p.OriginalFilename,
		ContentType:      p.ContentType,
		SizeBytes:        p.SizeBytes,
		CreatedAt:        p.CreatedAt,
	}
}
