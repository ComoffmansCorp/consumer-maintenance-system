package consumer

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/consumer/db"
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

func (r *Repository) Create(ctx context.Context, tenantID int64, name string, typ Type, description string) (Consumer, error) {
	row, err := r.queriesFor(ctx).CreateConsumer(ctx, db.CreateConsumerParams{
		Name:        strings.TrimSpace(name),
		Type:        string(typ),
		Description: textOrNull(description),
		TenantID:    tenantID,
	})
	if err != nil {
		return Consumer{}, fmt.Errorf("create consumer: %w", err)
	}
	return toConsumer(row), nil
}

func (r *Repository) GetByID(ctx context.Context, id, tenantID int64) (Consumer, error) {
	row, err := r.queriesFor(ctx).GetConsumerByID(ctx, db.GetConsumerByIDParams{ID: id, TenantID: tenantID})
	if err != nil {
		return Consumer{}, fmt.Errorf("get consumer by id: %w", err)
	}
	return toConsumer(row), nil
}

func (r *Repository) Exists(ctx context.Context, id, tenantID int64) (bool, error) {
	exists, err := r.queriesFor(ctx).ExistsConsumerByID(ctx, db.ExistsConsumerByIDParams{ID: id, TenantID: tenantID})
	if err != nil {
		return false, fmt.Errorf("exists consumer by id: %w", err)
	}
	return exists, nil
}

func (r *Repository) List(ctx context.Context, tenantID int64, search string, limit, offset int32) ([]Consumer, error) {
	rows, err := r.queriesFor(ctx).ListConsumers(ctx, db.ListConsumersParams{
		TenantID:   tenantID,
		Search:     strings.TrimSpace(search),
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list consumers: %w", err)
	}
	out := make([]Consumer, 0, len(rows))
	for _, row := range rows {
		out = append(out, toConsumer(row))
	}
	return out, nil
}

func (r *Repository) Count(ctx context.Context, tenantID int64, search string) (int64, error) {
	count, err := r.queriesFor(ctx).CountConsumers(ctx, db.CountConsumersParams{
		TenantID: tenantID,
		Search:   strings.TrimSpace(search),
	})
	if err != nil {
		return 0, fmt.Errorf("count consumers: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id, tenantID int64, name string, typ Type, description string) (Consumer, error) {
	row, err := r.queriesFor(ctx).UpdateConsumer(ctx, db.UpdateConsumerParams{
		Name:        strings.TrimSpace(name),
		Type:        string(typ),
		Description: textOrNull(description),
		ID:          id,
		TenantID:    tenantID,
	})
	if err != nil {
		return Consumer{}, fmt.Errorf("update consumer: %w", err)
	}
	return toConsumer(row), nil
}

func (r *Repository) Deactivate(ctx context.Context, id, tenantID int64) (Consumer, error) {
	row, err := r.queriesFor(ctx).DeactivateConsumer(ctx, db.DeactivateConsumerParams{ID: id, TenantID: tenantID})
	if err != nil {
		return Consumer{}, fmt.Errorf("deactivate consumer: %w", err)
	}
	return toConsumer(row), nil
}

func textOrNull(s string) pgtype.Text {
	s = strings.TrimSpace(s)
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func toConsumer(row db.Organization) Consumer {
	return Consumer{
		ID:          row.ID,
		Name:        row.Name,
		Type:        Type(row.Type),
		Description: row.Description.String,
		TenantID:    row.TenantID,
		Active:      row.Active,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
