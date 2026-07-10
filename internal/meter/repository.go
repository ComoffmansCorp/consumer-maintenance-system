package meter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/meter/db"
	platformdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/db"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
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

func (r *Repository) Create(ctx context.Context, actID int64, req CreateRequest) (Meter, error) {
	row, err := r.queriesFor(ctx).CreateMeter(ctx, db.CreateMeterParams{
		Type:                string(req.Type),
		SerialNumber:        strings.TrimSpace(req.SerialNumber),
		ManufactureYear:     int4OrNull(req.ManufactureYear),
		VerificationDate:    dateOrNull(httpx.DateToTime(req.VerificationDate)),
		SealState:           textOrNull(string(req.SealState)),
		TransformationRatio: int4OrNull(req.TransformationRatio),
		InspectionActID:     pgtype.Int8{Int64: actID, Valid: true},
	})
	if err != nil {
		return Meter{}, fmt.Errorf("create meter: %w", err)
	}
	return toMeter(row), nil
}

func (r *Repository) GetByID(ctx context.Context, id, actID int64) (Meter, error) {
	row, err := r.queriesFor(ctx).GetMeterByID(ctx, db.GetMeterByIDParams{
		ID:              id,
		InspectionActID: pgtype.Int8{Int64: actID, Valid: true},
	})
	if err != nil {
		return Meter{}, fmt.Errorf("get meter by id: %w", err)
	}
	return toMeter(row), nil
}

func (r *Repository) ListByAct(ctx context.Context, actID int64) ([]Meter, error) {
	rows, err := r.queriesFor(ctx).ListMetersByAct(ctx, pgtype.Int8{Int64: actID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("list meters by act: %w", err)
	}
	out := make([]Meter, 0, len(rows))
	for _, row := range rows {
		out = append(out, toMeter(row))
	}
	return out, nil
}

func (r *Repository) Update(ctx context.Context, id, actID int64, req UpdateRequest) (Meter, error) {
	row, err := r.queriesFor(ctx).UpdateMeter(ctx, db.UpdateMeterParams{
		Type:                string(req.Type),
		SerialNumber:        strings.TrimSpace(req.SerialNumber),
		ManufactureYear:     int4OrNull(req.ManufactureYear),
		VerificationDate:    dateOrNull(httpx.DateToTime(req.VerificationDate)),
		SealState:           textOrNull(string(req.SealState)),
		TransformationRatio: int4OrNull(req.TransformationRatio),
		ID:                  id,
		InspectionActID:     pgtype.Int8{Int64: actID, Valid: true},
	})
	if err != nil {
		return Meter{}, fmt.Errorf("update meter: %w", err)
	}
	return toMeter(row), nil
}

func (r *Repository) Delete(ctx context.Context, id, actID int64) error {
	if err := r.queriesFor(ctx).DeleteMeter(ctx, db.DeleteMeterParams{
		ID:              id,
		InspectionActID: pgtype.Int8{Int64: actID, Valid: true},
	}); err != nil {
		return fmt.Errorf("delete meter: %w", err)
	}
	return nil
}

func int4OrNull(v *int32) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *v, Valid: true}
}

func dateOrNull(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func textOrNull(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func toMeter(row db.Meter) Meter {
	var year *int32
	if row.ManufactureYear.Valid {
		year = &row.ManufactureYear.Int32
	}
	var verification *time.Time
	if row.VerificationDate.Valid {
		verification = &row.VerificationDate.Time
	}
	var ratio *int32
	if row.TransformationRatio.Valid {
		ratio = &row.TransformationRatio.Int32
	}
	return Meter{
		ID:                  row.ID,
		Type:                Type(row.Type),
		SerialNumber:        row.SerialNumber,
		ManufactureYear:     year,
		VerificationDate:    verification,
		SealState:           SealState(row.SealState.String),
		TransformationRatio: ratio,
		InspectionActID:     row.InspectionActID.Int64,
		CreatedAt:           row.CreatedAt.Time,
	}
}
