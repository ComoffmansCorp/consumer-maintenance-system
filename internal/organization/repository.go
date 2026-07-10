package organization

import (
	"context"
	"fmt"
	"strings"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/organization/db"
	platformdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/db"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{queries: queries}
}

func (r *Repository) queriesFor(ctx context.Context) *db.Queries {
	if tx, ok := platformdb.TxFromContext(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *Repository) GetByID(ctx context.Context, id int64) (Tenant, error) {
	row, err := r.queriesFor(ctx).GetTenantByID(ctx, id)
	if err != nil {
		return Tenant{}, fmt.Errorf("get tenant by id: %w", err)
	}
	return toTenant(row), nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (Tenant, error) {
	row, err := r.queriesFor(ctx).GetTenantByCode(ctx, strings.ToLower(strings.TrimSpace(code)))
	if err != nil {
		return Tenant{}, fmt.Errorf("get tenant by code: %w", err)
	}
	return toTenant(row), nil
}

func (r *Repository) ExistsByName(ctx context.Context, name string) (bool, error) {
	exists, err := r.queriesFor(ctx).ExistsTenantByName(ctx, strings.TrimSpace(name))
	if err != nil {
		return false, fmt.Errorf("exists tenant by name: %w", err)
	}
	return exists, nil
}

func (r *Repository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	exists, err := r.queriesFor(ctx).ExistsTenantByCode(ctx, NormalizeTenantCode(code))
	if err != nil {
		return false, fmt.Errorf("exists tenant by code: %w", err)
	}
	return exists, nil
}

func (r *Repository) Create(ctx context.Context, name, code string, plan TenantPlan) (Tenant, error) {
	row, err := r.queriesFor(ctx).CreateTenant(ctx, db.CreateTenantParams{
		Name:   strings.TrimSpace(name),
		Code:   NormalizeTenantCode(code),
		Plan:   string(plan),
		Active: true,
	})
	if err != nil {
		return Tenant{}, fmt.Errorf("create tenant: %w", err)
	}
	return toTenant(row), nil
}

func (r *Repository) List(ctx context.Context) ([]Tenant, error) {
	rows, err := r.queriesFor(ctx).ListTenants(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tenants: %w", err)
	}
	out := make([]Tenant, 0, len(rows))
	for _, row := range rows {
		out = append(out, toTenant(row))
	}
	return out, nil
}

func (r *Repository) UpdatePlan(ctx context.Context, id int64, plan TenantPlan) (Tenant, error) {
	row, err := r.queriesFor(ctx).UpdateTenantPlan(ctx, db.UpdateTenantPlanParams{
		ID:   id,
		Plan: string(plan),
	})
	if err != nil {
		return Tenant{}, fmt.Errorf("update tenant plan: %w", err)
	}
	return toTenant(row), nil
}

func (r *Repository) Deactivate(ctx context.Context, id int64) (Tenant, error) {
	row, err := r.queriesFor(ctx).DeactivateTenant(ctx, id)
	if err != nil {
		return Tenant{}, fmt.Errorf("deactivate tenant: %w", err)
	}
	return toTenant(row), nil
}

func NormalizeTenantCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

func toTenant(row db.Tenant) Tenant {
	return Tenant{
		ID:        row.ID,
		Name:      row.Name,
		Code:      row.Code,
		Plan:      TenantPlan(row.Plan),
		Active:    row.Active,
		CreatedAt: row.CreatedAt.Time,
	}
}

func ToDTO(t Tenant) TenantDTO {
	return TenantDTO{
		ID:     t.ID,
		Name:   t.Name,
		Code:   t.Code,
		Plan:   t.Plan,
		Active: t.Active,
	}
}
