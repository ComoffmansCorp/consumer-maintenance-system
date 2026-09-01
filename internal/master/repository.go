package master

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/master/db"
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

func (r *Repository) UpsertProfile(ctx context.Context, userID int64, city, bio string) (Profile, error) {
	row, err := r.queriesFor(ctx).UpsertMasterProfile(ctx, db.UpsertMasterProfileParams{
		UserID: userID,
		City:   pgtype.Text{String: city, Valid: city != ""},
		Bio:    pgtype.Text{String: bio, Valid: bio != ""},
	})
	if err != nil {
		return Profile{}, fmt.Errorf("upsert master profile: %w", err)
	}
	return toProfile(row), nil
}

func (r *Repository) GetProfile(ctx context.Context, userID int64) (Profile, error) {
	row, err := r.queriesFor(ctx).GetMasterProfile(ctx, userID)
	if err != nil {
		return Profile{}, fmt.Errorf("get master profile: %w", err)
	}
	return toProfile(row), nil
}

func (r *Repository) ListProfiles(ctx context.Context, limit, offset int32) ([]Profile, error) {
	rows, err := r.queriesFor(ctx).ListMasterProfiles(ctx, db.ListMasterProfilesParams{PageLimit: limit, PageOffset: offset})
	if err != nil {
		return nil, fmt.Errorf("list master profiles: %w", err)
	}
	out := make([]Profile, 0, len(rows))
	for _, row := range rows {
		out = append(out, toProfile(row))
	}
	return out, nil
}

func (r *Repository) CountProfiles(ctx context.Context) (int64, error) {
	count, err := r.queriesFor(ctx).CountMasterProfiles(ctx)
	if err != nil {
		return 0, fmt.Errorf("count master profiles: %w", err)
	}
	return count, nil
}

// RecordReview atomically folds a new rating into the running average.
func (r *Repository) RecordReview(ctx context.Context, masterUserID int64, rating int) (Profile, error) {
	row, err := r.queriesFor(ctx).RecordMasterReview(ctx, db.RecordMasterReviewParams{
		Rating: float64(rating),
		UserID: masterUserID,
	})
	if err != nil {
		return Profile{}, fmt.Errorf("record master review: %w", err)
	}
	return toProfile(row), nil
}

func (r *Repository) ReplaceSpecializations(ctx context.Context, masterUserID int64, serviceIDs []int64) error {
	q := r.queriesFor(ctx)
	if err := q.DeleteMasterSpecializations(ctx, masterUserID); err != nil {
		return fmt.Errorf("clear master specializations: %w", err)
	}
	for _, serviceID := range serviceIDs {
		if err := q.AddMasterSpecialization(ctx, db.AddMasterSpecializationParams{
			MasterUserID: masterUserID,
			ServiceID:    serviceID,
		}); err != nil {
			return fmt.Errorf("add master specialization: %w", err)
		}
	}
	return nil
}

func (r *Repository) ListSpecializationIDs(ctx context.Context, masterUserID int64) ([]int64, error) {
	ids, err := r.queriesFor(ctx).ListMasterSpecializationIDs(ctx, masterUserID)
	if err != nil {
		return nil, fmt.Errorf("list master specialization ids: %w", err)
	}
	return ids, nil
}

func (r *Repository) HasSpecialization(ctx context.Context, masterUserID, serviceID int64) (bool, error) {
	has, err := r.queriesFor(ctx).MasterHasSpecialization(ctx, db.MasterHasSpecializationParams{
		MasterUserID: masterUserID,
		ServiceID:    serviceID,
	})
	if err != nil {
		return false, fmt.Errorf("check master specialization: %w", err)
	}
	return has, nil
}

func toProfile(row db.MasterProfile) Profile {
	city, bio := "", ""
	if row.City.Valid {
		city = row.City.String
	}
	if row.Bio.Valid {
		bio = row.Bio.String
	}
	return Profile{
		UserID:      row.UserID,
		City:        city,
		Bio:         bio,
		RatingAvg:   row.RatingAvg,
		RatingCount: row.RatingCount,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
