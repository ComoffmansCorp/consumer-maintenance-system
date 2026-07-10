package notification

import "context"

// AdminLister is satisfied by the auth domain: it lets notification reach
// every tenant admin when a task reaches a terminal state, without owning
// user data itself.
type AdminLister interface {
	ListTenantAdminIDs(ctx context.Context, tenantID int64) ([]int64, error)
}
