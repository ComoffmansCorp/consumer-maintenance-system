package master

import "context"

// CatalogPort lets master validate specialization service IDs against the
// catalog domain without importing it directly. Satisfied by catalog.Service
// via a wiring adapter.
type CatalogPort interface {
	IsServiceActive(ctx context.Context, serviceID int64) (bool, error)
}
