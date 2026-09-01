package payment

import "time"

type PaymentDTO struct {
	ID          int64     `json:"id"`
	RequestID   int64     `json:"requestId"`
	Amount      float64   `json:"amount"`
	PlatformFee float64   `json:"platformFee"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func ToPaymentDTO(p Payment) PaymentDTO {
	return PaymentDTO(p)
}
