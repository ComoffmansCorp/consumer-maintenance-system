package task

import (
	"time"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

type DTO struct {
	ID           int64       `json:"id"`
	Type         Type        `json:"type"`
	AddressID    int64       `json:"addressId"`
	AddressLabel string      `json:"addressLabel,omitempty"`
	Status       Status      `json:"status"`
	DueDate      *httpx.Date `json:"dueDate,omitempty"`
	AssigneeID   *int64      `json:"assigneeId,omitempty"`
	AssigneeName string      `json:"assigneeName,omitempty"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	CompletedAt  *time.Time  `json:"completedAt,omitempty"`
	CanceledAt   *time.Time  `json:"canceledAt,omitempty"`
	CancelReason string      `json:"cancelReason,omitempty"`
}

type CreateRequest struct {
	Type       Type        `json:"type"`
	AddressID  int64       `json:"addressId"`
	DueDate    *httpx.Date `json:"dueDate"`
	AssigneeID *int64      `json:"assigneeId"`
}

type AssignRequest struct {
	AssigneeID int64 `json:"assigneeId"`
}

type StatusRequest struct {
	Status Status `json:"status"`
	Reason string `json:"reason"`
}

func ToDTO(t Task) DTO {
	return DTO{
		ID:           t.ID,
		Type:         t.Type,
		AddressID:    t.AddressID,
		Status:       t.Status,
		DueDate:      httpx.TimeToDate(t.DueDate),
		AssigneeID:   t.AssigneeID,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		CompletedAt:  t.CompletedAt,
		CanceledAt:   t.CanceledAt,
		CancelReason: t.CancelReason,
	}
}
