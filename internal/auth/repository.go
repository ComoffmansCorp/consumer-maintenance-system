package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/auth/db"
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

func (r *Repository) GetByUsername(ctx context.Context, username string) (User, error) {
	row, err := r.queriesFor(ctx).GetUserByUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		return User{}, fmt.Errorf("get user by username: %w", err)
	}
	return toUser(row), nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (User, error) {
	row, err := r.queriesFor(ctx).GetUserByID(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("get user by id: %w", err)
	}
	return toUser(row), nil
}

func (r *Repository) CountByRole(ctx context.Context, role Role) (int64, error) {
	count, err := r.queriesFor(ctx).CountUsersByRole(ctx, string(role))
	if err != nil {
		return 0, fmt.Errorf("count users by role: %w", err)
	}
	return count, nil
}

func (r *Repository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	exists, err := r.queriesFor(ctx).ExistsUserByUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		return false, fmt.Errorf("exists user by username: %w", err)
	}
	return exists, nil
}

func (r *Repository) Create(ctx context.Context, username, passwordHash, fullName string, role Role) (User, error) {
	row, err := r.queriesFor(ctx).CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		PasswordHash: passwordHash,
		FullName:     pgtype.Text{String: strings.TrimSpace(fullName), Valid: strings.TrimSpace(fullName) != ""},
		Role:         string(role),
	})
	if err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return toUser(row), nil
}

func (r *Repository) CreateRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	_, err := r.queriesFor(ctx).CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt.UTC(), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}
	return nil
}

func (r *Repository) GetRefreshToken(ctx context.Context, tokenHash string) (RefreshToken, error) {
	row, err := r.queriesFor(ctx).GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return RefreshToken{}, fmt.Errorf("get refresh token: %w", err)
	}
	return RefreshToken{
		ID:        row.ID,
		UserID:    row.UserID,
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt.Time,
	}, nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	if err := r.queriesFor(ctx).RevokeRefreshToken(ctx, tokenHash); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func toUser(row db.User) User {
	fullName := ""
	if row.FullName.Valid {
		fullName = row.FullName.String
	}
	return User{
		ID:           row.ID,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		FullName:     fullName,
		Role:         Role(row.Role),
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}
