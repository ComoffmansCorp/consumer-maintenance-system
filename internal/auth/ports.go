package auth

import "context"

type TenantInfo struct {
	ID        int64
	Name      string
	Code      string
	Plan      string
	Active    bool
	UserLimit int
}

type TenantStore interface {
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	Create(ctx context.Context, name, code, plan string) (TenantInfo, error)
	GetByCode(ctx context.Context, code string) (TenantInfo, error)
	GetByID(ctx context.Context, id int64) (TenantInfo, error)
}

type TxRunner interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
