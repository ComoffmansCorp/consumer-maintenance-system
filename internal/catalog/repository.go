package catalog

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/catalog/db"
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

func (r *Repository) CreateCategory(ctx context.Context, parentCategoryID *int64, name string) (Category, error) {
	row, err := r.queriesFor(ctx).CreateCategory(ctx, db.CreateCategoryParams{
		ParentCategoryID: int8OrNil(parentCategoryID),
		Name:             name,
	})
	if err != nil {
		return Category{}, fmt.Errorf("create category: %w", err)
	}
	return toCategory(row), nil
}

func (r *Repository) GetCategoryByID(ctx context.Context, id int64) (Category, error) {
	row, err := r.queriesFor(ctx).GetCategoryByID(ctx, id)
	if err != nil {
		return Category{}, fmt.Errorf("get category by id: %w", err)
	}
	return toCategory(row), nil
}

func (r *Repository) UpdateCategory(ctx context.Context, id int64, name string, active bool) (Category, error) {
	row, err := r.queriesFor(ctx).UpdateCategory(ctx, db.UpdateCategoryParams{ID: id, Name: name, Active: active})
	if err != nil {
		return Category{}, fmt.Errorf("update category: %w", err)
	}
	return toCategory(row), nil
}

func (r *Repository) ListActiveCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.queriesFor(ctx).ListActiveCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("list active categories: %w", err)
	}
	out := make([]Category, 0, len(rows))
	for _, row := range rows {
		out = append(out, toCategory(row))
	}
	return out, nil
}

func (r *Repository) CreateService(ctx context.Context, categoryID int64, name, description string, priceFrom, priceTo *float64, unit string) (Offering, error) {
	row, err := r.queriesFor(ctx).CreateService(ctx, db.CreateServiceParams{
		CategoryID:  categoryID,
		Name:        name,
		Description: pgtype.Text{String: description, Valid: description != ""},
		PriceFrom:   priceFrom,
		PriceTo:     priceTo,
		Unit:        pgtype.Text{String: unit, Valid: unit != ""},
	})
	if err != nil {
		return Offering{}, fmt.Errorf("create service: %w", err)
	}
	return toOffering(row), nil
}

func (r *Repository) GetServiceByID(ctx context.Context, id int64) (Offering, error) {
	row, err := r.queriesFor(ctx).GetServiceByID(ctx, id)
	if err != nil {
		return Offering{}, fmt.Errorf("get service by id: %w", err)
	}
	return toOffering(row), nil
}

func (r *Repository) UpdateService(ctx context.Context, id int64, name, description string, priceFrom, priceTo *float64, unit string, active bool) (Offering, error) {
	row, err := r.queriesFor(ctx).UpdateService(ctx, db.UpdateServiceParams{
		ID:          id,
		Name:        name,
		Description: pgtype.Text{String: description, Valid: description != ""},
		PriceFrom:   priceFrom,
		PriceTo:     priceTo,
		Unit:        pgtype.Text{String: unit, Valid: unit != ""},
		Active:      active,
	})
	if err != nil {
		return Offering{}, fmt.Errorf("update service: %w", err)
	}
	return toOffering(row), nil
}

func (r *Repository) ListActiveServices(ctx context.Context, categoryID int64) ([]Offering, error) {
	rows, err := r.queriesFor(ctx).ListActiveServices(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("list active services: %w", err)
	}
	out := make([]Offering, 0, len(rows))
	for _, row := range rows {
		out = append(out, toOffering(row))
	}
	return out, nil
}

func toCategory(row db.ServiceCategory) Category {
	var parentID *int64
	if row.ParentCategoryID.Valid {
		parentID = &row.ParentCategoryID.Int64
	}
	return Category{
		ID:               row.ID,
		ParentCategoryID: parentID,
		Name:             row.Name,
		Active:           row.Active,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}
}

func toOffering(row db.Service) Offering {
	description, unit := "", ""
	if row.Description.Valid {
		description = row.Description.String
	}
	if row.Unit.Valid {
		unit = row.Unit.String
	}
	return Offering{
		ID:          row.ID,
		CategoryID:  row.CategoryID,
		Name:        row.Name,
		Description: description,
		PriceFrom:   row.PriceFrom,
		PriceTo:     row.PriceTo,
		Unit:        unit,
		Active:      row.Active,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}

func int8OrNil(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}
