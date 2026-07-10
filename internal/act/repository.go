package act

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/act/db"
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

// --- Inspection acts ---

func (r *Repository) CreateInspection(ctx context.Context, taskID, tenantID, addressID int64, consumerID *int64, req CreateInspectionRequest) (InspectionAct, error) {
	row, err := r.queriesFor(ctx).CreateInspectionAct(ctx, db.CreateInspectionActParams{
		TaskID:         pgtype.Int8{Int64: taskID, Valid: true},
		TenantID:       tenantID,
		AddressID:      addressID,
		InspectionDate: dateOrNull(httpx.DateToTime(req.InspectionDate)),
		ConsumerID:     int8OrNull(consumerID),
		InspectionType: textOrNull(string(req.InspectionType)),
		Notes:          textOrNull(req.Notes),
	})
	if err != nil {
		return InspectionAct{}, fmt.Errorf("create inspection act: %w", err)
	}
	return toInspectionAct(row), nil
}

func (r *Repository) GetInspectionByID(ctx context.Context, id, tenantID int64) (InspectionAct, error) {
	row, err := r.queriesFor(ctx).GetInspectionActByID(ctx, db.GetInspectionActByIDParams{ID: id, TenantID: tenantID})
	if err != nil {
		return InspectionAct{}, fmt.Errorf("get inspection act by id: %w", err)
	}
	return toInspectionAct(row), nil
}

func (r *Repository) GetInspectionByTaskID(ctx context.Context, taskID, tenantID int64) (InspectionAct, error) {
	row, err := r.queriesFor(ctx).GetInspectionActByTaskID(ctx, db.GetInspectionActByTaskIDParams{
		TaskID:   pgtype.Int8{Int64: taskID, Valid: true},
		TenantID: tenantID,
	})
	if err != nil {
		return InspectionAct{}, fmt.Errorf("get inspection act by task id: %w", err)
	}
	return toInspectionAct(row), nil
}

func (r *Repository) ExistsInspectionByTaskID(ctx context.Context, taskID int64) (bool, error) {
	exists, err := r.queriesFor(ctx).ExistsInspectionActByTaskID(ctx, pgtype.Int8{Int64: taskID, Valid: true})
	if err != nil {
		return false, fmt.Errorf("exists inspection act by task id: %w", err)
	}
	return exists, nil
}

func (r *Repository) UpdateInspection(ctx context.Context, id, tenantID int64, consumerID *int64, req UpdateInspectionRequest) (InspectionAct, error) {
	row, err := r.queriesFor(ctx).UpdateInspectionAct(ctx, db.UpdateInspectionActParams{
		InspectionDate: dateOrNull(httpx.DateToTime(req.InspectionDate)),
		ConsumerID:     int8OrNull(consumerID),
		InspectionType: textOrNull(string(req.InspectionType)),
		Notes:          textOrNull(req.Notes),
		ID:             id,
		TenantID:       tenantID,
	})
	if err != nil {
		return InspectionAct{}, fmt.Errorf("update inspection act: %w", err)
	}
	return toInspectionAct(row), nil
}

// --- Replacement acts ---

func (r *Repository) CreateReplacement(ctx context.Context, taskID, tenantID, addressID int64, req CreateReplacementRequest) (ReplacementAct, error) {
	row, err := r.queriesFor(ctx).CreateReplacementAct(ctx, db.CreateReplacementActParams{
		TaskID:           pgtype.Int8{Int64: taskID, Valid: true},
		TenantID:         tenantID,
		AddressID:        addressID,
		AccountNumber:    strings.TrimSpace(req.AccountNumber),
		InstallationDate: dateOrNull(httpx.DateToTime(req.InstallationDate)),
		OldBrand:         textOrNull(req.OldBrand),
		OldSerialNumber:  textOrNull(req.OldSerialNumber),
		OldReadings:      float8OrNull(req.OldReadings),
		NewBrand:         textOrNull(req.NewBrand),
		NewSerialNumber:  textOrNull(req.NewSerialNumber),
		NewReadings:      float8OrNull(req.NewReadings),
	})
	if err != nil {
		return ReplacementAct{}, fmt.Errorf("create replacement act: %w", err)
	}
	return toReplacementAct(row), nil
}

func (r *Repository) GetReplacementByID(ctx context.Context, id, tenantID int64) (ReplacementAct, error) {
	row, err := r.queriesFor(ctx).GetReplacementActByID(ctx, db.GetReplacementActByIDParams{ID: id, TenantID: tenantID})
	if err != nil {
		return ReplacementAct{}, fmt.Errorf("get replacement act by id: %w", err)
	}
	return toReplacementAct(row), nil
}

func (r *Repository) GetReplacementByTaskID(ctx context.Context, taskID, tenantID int64) (ReplacementAct, error) {
	row, err := r.queriesFor(ctx).GetReplacementActByTaskID(ctx, db.GetReplacementActByTaskIDParams{
		TaskID:   pgtype.Int8{Int64: taskID, Valid: true},
		TenantID: tenantID,
	})
	if err != nil {
		return ReplacementAct{}, fmt.Errorf("get replacement act by task id: %w", err)
	}
	return toReplacementAct(row), nil
}

func (r *Repository) ExistsReplacementByTaskID(ctx context.Context, taskID int64) (bool, error) {
	exists, err := r.queriesFor(ctx).ExistsReplacementActByTaskID(ctx, pgtype.Int8{Int64: taskID, Valid: true})
	if err != nil {
		return false, fmt.Errorf("exists replacement act by task id: %w", err)
	}
	return exists, nil
}

func (r *Repository) UpdateReplacement(ctx context.Context, id, tenantID int64, req UpdateReplacementRequest) (ReplacementAct, error) {
	row, err := r.queriesFor(ctx).UpdateReplacementAct(ctx, db.UpdateReplacementActParams{
		AccountNumber:    strings.TrimSpace(req.AccountNumber),
		InstallationDate: dateOrNull(httpx.DateToTime(req.InstallationDate)),
		OldBrand:         textOrNull(req.OldBrand),
		OldSerialNumber:  textOrNull(req.OldSerialNumber),
		OldReadings:      float8OrNull(req.OldReadings),
		NewBrand:         textOrNull(req.NewBrand),
		NewSerialNumber:  textOrNull(req.NewSerialNumber),
		NewReadings:      float8OrNull(req.NewReadings),
		ID:               id,
		TenantID:         tenantID,
	})
	if err != nil {
		return ReplacementAct{}, fmt.Errorf("update replacement act: %w", err)
	}
	return toReplacementAct(row), nil
}

// --- helpers ---

func dateOrNull(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func int8OrNull(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}

func float8OrNull(v *float64) pgtype.Float8 {
	if v == nil {
		return pgtype.Float8{}
	}
	return pgtype.Float8{Float64: *v, Valid: true}
}

func textOrNull(s string) pgtype.Text {
	s = strings.TrimSpace(s)
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func toInspectionAct(row db.InspectionAct) InspectionAct {
	var inspectionDate *time.Time
	if row.InspectionDate.Valid {
		inspectionDate = &row.InspectionDate.Time
	}
	var consumerID *int64
	if row.ConsumerID.Valid {
		consumerID = &row.ConsumerID.Int64
	}
	return InspectionAct{
		ID:             row.ID,
		TaskID:         row.TaskID.Int64,
		TenantID:       row.TenantID,
		AddressID:      row.AddressID,
		InspectionDate: inspectionDate,
		ConsumerID:     consumerID,
		InspectionType: InspectionType(row.InspectionType.String),
		Notes:          row.Notes.String,
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
	}
}

func toReplacementAct(row db.ReplacementAct) ReplacementAct {
	var installationDate *time.Time
	if row.InstallationDate.Valid {
		installationDate = &row.InstallationDate.Time
	}
	var oldReadings *float64
	if row.OldReadings.Valid {
		oldReadings = &row.OldReadings.Float64
	}
	var newReadings *float64
	if row.NewReadings.Valid {
		newReadings = &row.NewReadings.Float64
	}
	return ReplacementAct{
		ID:               row.ID,
		TaskID:           row.TaskID.Int64,
		TenantID:         row.TenantID,
		AddressID:        row.AddressID,
		AccountNumber:    row.AccountNumber,
		InstallationDate: installationDate,
		OldBrand:         row.OldBrand.String,
		OldSerialNumber:  row.OldSerialNumber.String,
		OldReadings:      oldReadings,
		NewBrand:         row.NewBrand.String,
		NewSerialNumber:  row.NewSerialNumber.String,
		NewReadings:      newReadings,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}
}
