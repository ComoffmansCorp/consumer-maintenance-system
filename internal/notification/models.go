package notification

import "time"

const (
	TypeTaskAssigned      = "TASK_ASSIGNED"
	TypeTaskStatusChanged = "TASK_STATUS_CHANGED"
)

type Notification struct {
	ID        int64
	TenantID  int64
	UserID    int64
	Type      string
	Title     string
	Message   string
	Payload   map[string]any
	ReadAt    *time.Time
	CreatedAt time.Time
}
