package task

import "time"

type Type string

const (
	TypeInspection  Type = "INSPECTION"
	TypeReplacement Type = "REPLACEMENT"
)

func (t Type) Valid() bool {
	switch t {
	case TypeInspection, TypeReplacement:
		return true
	default:
		return false
	}
}

type Status string

const (
	StatusPending    Status = "PENDING"
	StatusInProgress Status = "IN_PROGRESS"
	StatusCompleted  Status = "COMPLETED"
	StatusCanceled   Status = "CANCELED"
)

type Task struct {
	ID           int64
	Type         Type
	TenantID     int64
	AddressID    int64
	Status       Status
	DueDate      *time.Time
	AssigneeID   *int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CompletedAt  *time.Time
	CanceledAt   *time.Time
	CancelReason string
}
