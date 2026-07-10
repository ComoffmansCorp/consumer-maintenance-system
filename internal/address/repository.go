package address

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/address/db"
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

func (r *Repository) Create(ctx context.Context, tenantID int64, street, house, building, apartment string, consumerID *int64) (Address, error) {
	row, err := r.queriesFor(ctx).CreateAddress(ctx, db.CreateAddressParams{
		Street:     strings.TrimSpace(street),
		House:      strings.TrimSpace(house),
		Building:   textOrNull(building),
		Apartment:  textOrNull(apartment),
		TenantID:   tenantID,
		ConsumerID: int8OrNull(consumerID),
	})
	if err != nil {
		return Address{}, fmt.Errorf("create address: %w", err)
	}
	return toAddress(row), nil
}

func (r *Repository) GetByID(ctx context.Context, id, tenantID int64) (Address, error) {
	row, err := r.queriesFor(ctx).GetAddressByID(ctx, db.GetAddressByIDParams{ID: id, TenantID: tenantID})
	if err != nil {
		return Address{}, fmt.Errorf("get address by id: %w", err)
	}
	return toAddress(row), nil
}

func (r *Repository) Exists(ctx context.Context, id, tenantID int64) (bool, error) {
	exists, err := r.queriesFor(ctx).ExistsAddressByID(ctx, db.ExistsAddressByIDParams{ID: id, TenantID: tenantID})
	if err != nil {
		return false, fmt.Errorf("exists address by id: %w", err)
	}
	return exists, nil
}

func (r *Repository) List(ctx context.Context, tenantID int64, search string, consumerID int64, limit, offset int32) ([]Address, error) {
	rows, err := r.queriesFor(ctx).ListAddresses(ctx, db.ListAddressesParams{
		TenantID:   tenantID,
		Search:     strings.TrimSpace(search),
		ConsumerID: consumerID,
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list addresses: %w", err)
	}
	out := make([]Address, 0, len(rows))
	for _, row := range rows {
		out = append(out, toAddress(row))
	}
	return out, nil
}

func (r *Repository) Count(ctx context.Context, tenantID int64, search string, consumerID int64) (int64, error) {
	count, err := r.queriesFor(ctx).CountAddresses(ctx, db.CountAddressesParams{
		TenantID:   tenantID,
		Search:     strings.TrimSpace(search),
		ConsumerID: consumerID,
	})
	if err != nil {
		return 0, fmt.Errorf("count addresses: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id, tenantID int64, street, house, building, apartment string, consumerID *int64) (Address, error) {
	row, err := r.queriesFor(ctx).UpdateAddress(ctx, db.UpdateAddressParams{
		Street:     strings.TrimSpace(street),
		House:      strings.TrimSpace(house),
		Building:   textOrNull(building),
		Apartment:  textOrNull(apartment),
		ConsumerID: int8OrNull(consumerID),
		ID:         id,
		TenantID:   tenantID,
	})
	if err != nil {
		return Address{}, fmt.Errorf("update address: %w", err)
	}
	return toAddress(row), nil
}

func textOrNull(s string) pgtype.Text {
	s = strings.TrimSpace(s)
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func int8OrNull(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}

func toAddress(row db.Address) Address {
	var consumerID *int64
	if row.ConsumerID.Valid {
		consumerID = &row.ConsumerID.Int64
	}
	return Address{
		ID:         row.ID,
		Street:     row.Street,
		House:      row.House,
		Building:   row.Building.String,
		Apartment:  row.Apartment.String,
		TenantID:   row.TenantID,
		ConsumerID: consumerID,
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
	}
}
