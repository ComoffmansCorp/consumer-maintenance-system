package task

import "context"

type AddressStore interface {
	Exists(ctx context.Context, id, tenantID int64) (bool, error)
	GetLabel(ctx context.Context, id, tenantID int64) (string, error)
}

type AssigneeSummary struct {
	ID       int64
	FullName string
	Username string
}

// UserStore is satisfied by the auth domain: it lets task validate that a
// candidate assignee exists, belongs to the tenant, and holds the
// ELECTRICIAN role, without task owning any user data itself.
type UserStore interface {
	GetElectrician(ctx context.Context, id, tenantID int64) (AssigneeSummary, error)
}

// ActStore lets task confirm the corresponding act (inspection or
// replacement, matching the task type) was filled in before the task can be
// marked completed.
type ActStore interface {
	HasActForTask(ctx context.Context, taskID int64, taskType string) (bool, error)
}
