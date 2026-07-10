package wiring

import (
	"context"
	"fmt"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/auth"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/organization"
)

type TenantAdapter struct {
	org *organization.Repository
}

func NewTenantAdapter(org *organization.Repository) *TenantAdapter {
	return &TenantAdapter{org: org}
}

func (a *TenantAdapter) ExistsByName(ctx context.Context, name string) (bool, error) {
	return a.org.ExistsByName(ctx, name)
}

func (a *TenantAdapter) ExistsByCode(ctx context.Context, code string) (bool, error) {
	return a.org.ExistsByCode(ctx, code)
}

func (a *TenantAdapter) Create(ctx context.Context, name, code, plan string) (auth.TenantInfo, error) {
	tenant, err := a.org.Create(ctx, name, code, organization.TenantPlan(plan))
	if err != nil {
		return auth.TenantInfo{}, err
	}
	return toTenantInfo(tenant), nil
}

func (a *TenantAdapter) GetByCode(ctx context.Context, code string) (auth.TenantInfo, error) {
	tenant, err := a.org.GetByCode(ctx, code)
	if err != nil {
		return auth.TenantInfo{}, err
	}
	return toTenantInfo(tenant), nil
}

func (a *TenantAdapter) GetByID(ctx context.Context, id int64) (auth.TenantInfo, error) {
	tenant, err := a.org.GetByID(ctx, id)
	if err != nil {
		return auth.TenantInfo{}, err
	}
	return toTenantInfo(tenant), nil
}

func toTenantInfo(t organization.Tenant) auth.TenantInfo {
	return auth.TenantInfo{
		ID:        t.ID,
		Name:      t.Name,
		Code:      t.Code,
		Plan:      string(t.Plan),
		Active:    t.Active,
		UserLimit: t.Plan.UserLimit(),
	}
}

type TxRunnerAdapter struct {
	manager interface {
		WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	}
}

func NewTxRunnerAdapter(manager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}) *TxRunnerAdapter {
	return &TxRunnerAdapter{manager: manager}
}

func (a *TxRunnerAdapter) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := a.manager.WithinTransaction(ctx, fn); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}
