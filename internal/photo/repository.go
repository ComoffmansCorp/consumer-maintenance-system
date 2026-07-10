package photo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/photo/db"
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

type CreateParams struct {
	Filename         string
	Note             string
	TenantID         int64
	InspectionActID  *int64
	ReplacementActID *int64
	OriginalFilename string
	ContentType      string
	SizeBytes        int64
	UploadedBy       *int64
}

func (r *Repository) Create(ctx context.Context, p CreateParams) (Photo, error) {
	row, err := r.queriesFor(ctx).CreatePhoto(ctx, db.CreatePhotoParams{
		Filename:         p.Filename,
		Note:             textOrNull(p.Note),
		TenantID:         p.TenantID,
		InspectionActID:  int8OrNull(p.InspectionActID),
		ReplacementActID: int8OrNull(p.ReplacementActID),
		OriginalFilename: p.OriginalFilename,
		ContentType:      p.ContentType,
		SizeBytes:        p.SizeBytes,
		UploadedBy:       int8OrNull(p.UploadedBy),
	})
	if err != nil {
		return Photo{}, fmt.Errorf("create photo: %w", err)
	}
	return toPhoto(row), nil
}

func (r *Repository) GetByID(ctx context.Context, id, tenantID int64) (Photo, error) {
	row, err := r.queriesFor(ctx).GetPhotoByID(ctx, db.GetPhotoByIDParams{ID: id, TenantID: tenantID})
	if err != nil {
		return Photo{}, fmt.Errorf("get photo by id: %w", err)
	}
	return toPhoto(row), nil
}

func (r *Repository) ListByInspectionAct(ctx context.Context, actID int64) ([]Photo, error) {
	rows, err := r.queriesFor(ctx).ListPhotosByInspectionAct(ctx, pgtype.Int8{Int64: actID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("list photos by inspection act: %w", err)
	}
	out := make([]Photo, 0, len(rows))
	for _, row := range rows {
		out = append(out, toPhoto(row))
	}
	return out, nil
}

func (r *Repository) ListByReplacementAct(ctx context.Context, actID int64) ([]Photo, error) {
	rows, err := r.queriesFor(ctx).ListPhotosByReplacementAct(ctx, pgtype.Int8{Int64: actID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("list photos by replacement act: %w", err)
	}
	out := make([]Photo, 0, len(rows))
	for _, row := range rows {
		out = append(out, toPhoto(row))
	}
	return out, nil
}

func (r *Repository) Delete(ctx context.Context, id, tenantID int64) (Photo, error) {
	row, err := r.queriesFor(ctx).DeletePhoto(ctx, db.DeletePhotoParams{ID: id, TenantID: tenantID})
	if err != nil {
		return Photo{}, fmt.Errorf("delete photo: %w", err)
	}
	return toPhoto(row), nil
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

func toPhoto(row db.Photo) Photo {
	var inspectionActID *int64
	if row.InspectionActID.Valid {
		inspectionActID = &row.InspectionActID.Int64
	}
	var replacementActID *int64
	if row.ReplacementActID.Valid {
		replacementActID = &row.ReplacementActID.Int64
	}
	var uploadedBy *int64
	if row.UploadedBy.Valid {
		uploadedBy = &row.UploadedBy.Int64
	}
	return Photo{
		ID:               row.ID,
		Filename:         row.Filename,
		Note:             row.Note.String,
		TenantID:         row.TenantID,
		InspectionActID:  inspectionActID,
		ReplacementActID: replacementActID,
		OriginalFilename: row.OriginalFilename,
		ContentType:      row.ContentType,
		SizeBytes:        row.SizeBytes,
		UploadedBy:       uploadedBy,
		CreatedAt:        row.CreatedAt.Time,
	}
}
