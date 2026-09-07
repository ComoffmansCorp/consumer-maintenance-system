package request

import "context"

// ServiceInfo is the sliver of catalog data request needs: enough to
// validate a serviceId and to enrich a request DTO with a display name.
type ServiceInfo struct {
	ID     int64
	Name   string
	Active bool
}

// CatalogPort lets request validate/enrich the service a заявка references,
// without importing catalog directly. Satisfied by catalog.Service via a
// wiring adapter. found=false means the id doesn't exist at all.
type CatalogPort interface {
	GetService(ctx context.Context, serviceID int64) (info ServiceInfo, found bool, err error)
}

// SpecializationPort is the hard, server-side specialization gate on offer
// submission -- same rule the old flat marketplace package enforced on
// claim: a master cannot bid on work outside what they're specialized in,
// even by calling the API directly. Also doubles as the lookup ListOffers
// uses to enrich each OfferDTO with the bidding master's avatar, since it's
// the same master.Service already wired in as this port. Satisfied by
// master.Service via a wiring adapter.
type SpecializationPort interface {
	HasSpecialization(ctx context.Context, masterUserID, serviceID int64) (bool, error)
	GetAvatarURL(ctx context.Context, masterUserID int64) (*string, error)
}

// TxRunner wraps multi-step writes (AcceptOffer touches the offer, its
// siblings, and the request in one commit) in a single transaction.
// Satisfied by platformdb.TxManager via a wiring adapter.
type TxRunner interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
