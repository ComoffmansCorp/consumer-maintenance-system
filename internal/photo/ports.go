package photo

import "context"

// ActStore lets the photo domain confirm the target act exists and belongs
// to the current tenant before attaching or serving a file, without owning
// the acts tables itself.
type ActStore interface {
	EnsureInspectionAct(ctx context.Context, actID, tenantID int64) error
	EnsureReplacementAct(ctx context.Context, actID, tenantID int64) error
}
