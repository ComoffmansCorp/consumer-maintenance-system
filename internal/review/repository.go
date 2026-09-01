package review

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	platformdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/db"
	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/review/db"
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

func (r *Repository) Create(ctx context.Context, requestID, clientID, masterID int64, rating int, comment string) (Review, error) {
	row, err := r.queriesFor(ctx).CreateReview(ctx, db.CreateReviewParams{
		RequestID: requestID,
		ClientID:  clientID,
		MasterID:  masterID,
		Rating:    int16(rating),
		Comment:   pgtype.Text{String: comment, Valid: comment != ""},
	})
	if err != nil {
		return Review{}, fmt.Errorf("create review: %w", err)
	}
	return toReview(row), nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (Review, error) {
	row, err := r.queriesFor(ctx).GetReviewByID(ctx, id)
	if err != nil {
		return Review{}, fmt.Errorf("get review by id: %w", err)
	}
	return toReview(row), nil
}

func (r *Repository) ListVisibleByMaster(ctx context.Context, masterID int64, limit, offset int32) ([]Review, error) {
	rows, err := r.queriesFor(ctx).ListVisibleReviewsByMaster(ctx, db.ListVisibleReviewsByMasterParams{
		MasterID: masterID, PageLimit: limit, PageOffset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list visible reviews by master: %w", err)
	}
	out := make([]Review, 0, len(rows))
	for _, row := range rows {
		out = append(out, toReview(row))
	}
	return out, nil
}

func (r *Repository) CountVisibleByMaster(ctx context.Context, masterID int64) (int64, error) {
	count, err := r.queriesFor(ctx).CountVisibleReviewsByMaster(ctx, masterID)
	if err != nil {
		return 0, fmt.Errorf("count visible reviews by master: %w", err)
	}
	return count, nil
}

func (r *Repository) Hide(ctx context.Context, id int64) (Review, error) {
	row, err := r.queriesFor(ctx).HideReview(ctx, id)
	if err != nil {
		return Review{}, fmt.Errorf("hide review: %w", err)
	}
	return toReview(row), nil
}

func toReview(row db.Review) Review {
	comment := ""
	if row.Comment.Valid {
		comment = row.Comment.String
	}
	return Review{
		ID:        row.ID,
		RequestID: row.RequestID,
		ClientID:  row.ClientID,
		MasterID:  row.MasterID,
		Rating:    int(row.Rating),
		Comment:   comment,
		Hidden:    row.Hidden,
		CreatedAt: row.CreatedAt.Time,
	}
}
