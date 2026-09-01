package marketplace

import "context"

type UserSummary struct {
	ID       int64
	FullName string
}

// UserStore lets marketplace resolve display names for clients and masters
// without owning user data itself. Satisfied by the auth domain.
type UserStore interface {
	GetUser(ctx context.Context, id int64) (UserSummary, error)
}

type TxRunner interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
