package meter

import (
	"time"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

type DTO struct {
	ID                  int64       `json:"id"`
	Type                Type        `json:"type"`
	SerialNumber        string      `json:"serialNumber"`
	ManufactureYear     *int32      `json:"manufactureYear,omitempty"`
	VerificationDate    *httpx.Date `json:"verificationDate,omitempty"`
	SealState           SealState   `json:"sealState,omitempty"`
	TransformationRatio *int32      `json:"transformationRatio,omitempty"`
	CreatedAt           time.Time   `json:"createdAt"`
}

type CreateRequest struct {
	Type                Type        `json:"type"`
	SerialNumber        string      `json:"serialNumber"`
	ManufactureYear     *int32      `json:"manufactureYear"`
	VerificationDate    *httpx.Date `json:"verificationDate"`
	SealState           SealState   `json:"sealState"`
	TransformationRatio *int32      `json:"transformationRatio"`
}

type UpdateRequest = CreateRequest

func ToDTO(m Meter) DTO {
	return DTO{
		ID:                  m.ID,
		Type:                m.Type,
		SerialNumber:        m.SerialNumber,
		ManufactureYear:     m.ManufactureYear,
		VerificationDate:    httpx.TimeToDate(m.VerificationDate),
		SealState:           m.SealState,
		TransformationRatio: m.TransformationRatio,
		CreatedAt:           m.CreatedAt,
	}
}
