package meter

import "context"

// ActStore lets the meter domain confirm an inspection act exists and
// belongs to the current tenant before mutating its meters, without a
// cross-domain join or duplicating the acts table.
type ActStore interface {
	EnsureInspectionAct(ctx context.Context, actID, tenantID int64) error
}
