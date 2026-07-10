package address

import "context"

// ConsumerSummary is the read-only projection the address domain needs from
// the consumer domain to enrich responses without a cross-domain join.
type ConsumerSummary struct {
	ID   int64
	Name string
	Type string
}

type ConsumerStore interface {
	Exists(ctx context.Context, id, tenantID int64) (bool, error)
	GetSummary(ctx context.Context, id, tenantID int64) (ConsumerSummary, error)
}
