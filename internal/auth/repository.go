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

func (r *Repository) GetPlatformUser(ctx context.Context, username string) (User, error) {
	row, err := r.queriesFor(ctx).GetUserByPlatformUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		return User{}, fmt.Errorf("get platform user: %w", err)
	}
	return toUser(row), nil
}

func (r *Repository) GetTenantUser(ctx context.Context, tenantCode, username string) (User, error) {
	row, err := r.queriesFor(ctx).GetUserByTenantCodeAndUsername(ctx, db.GetUserByTenantCodeAndUsernameParams{
		TenantCode: strings.TrimSpace(tenantCode),
		Username:   strings.TrimSpace(username),
	})
	if err != nil {
		return User{}, fmt.Errorf("get tenant user: %w", err)
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

func (r *Repository) CountByTenantID(ctx context.Context, tenantID int64) (int64, error) {
	count, err := r.queriesFor(ctx).CountUsersByTenantID(ctx, tenantID)
	if err != nil {
		return 0, fmt.Errorf("count users by tenant: %w", err)
	}
	return count, nil
}

func (r *Repository) ListByTenant(ctx context.Context, tenantID int64, role Role, limit, offset int32) ([]User, error) {
	rows, err := r.queriesFor(ctx).ListUsersByTenant(ctx, db.ListUsersByTenantParams{
		TenantID:   tenantID,
		Role:       string(role),
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list users by tenant: %w", err)
	}
	out := make([]User, 0, len(rows))
	for _, row := range rows {
		out = append(out, toUser(row))
	}
	return out, nil
}

func (r *Repository) CountByTenantAndRole(ctx context.Context, tenantID int64, role Role) (int64, error) {
	count, err := r.queriesFor(ctx).CountUsersByTenantAndRole(ctx, db.CountUsersByTenantAndRoleParams{
		TenantID: tenantID,
		Role:     string(role),
	})
	if err != nil {
		return 0, fmt.Errorf("count users by tenant and role: %w", err)
	}
	return count, nil
}

func (r *Repository) ExistsInTenant(ctx context.Context, tenantID int64, username string) (bool, error) {
	exists, err := r.queriesFor(ctx).ExistsUserInTenant(ctx, db.ExistsUserInTenantParams{
		TenantID: pgtype.Int8{Int64: tenantID, Valid: true},
		Username: strings.TrimSpace(username),
	})
	if err != nil {
		return false, fmt.Errorf("exists user in tenant: %w", err)
	}
	return exists, nil
}

func (r *Repository) ExistsPlatformUser(ctx context.Context, username string) (bool, error) {
	exists, err := r.queriesFor(ctx).ExistsPlatformUser(ctx, strings.TrimSpace(username))
	if err != nil {
		return false, fmt.Errorf("exists platform user: %w", err)
	}
	return exists, nil
}

func (r *Repository) Create(ctx context.Context, username, passwordHash, fullName string, role Role, tenantID *int64) (User, error) {
	var tenant pgtype.Int8
	if tenantID != nil {
		tenant = pgtype.Int8{Int64: *tenantID, Valid: true}
	}

	row, err := r.queriesFor(ctx).CreateUser(ctx, db.CreateUserParams{
		Username: username,
		Password: passwordHash,
		FullName: pgtype.Text{String: strings.TrimSpace(fullName), Valid: strings.TrimSpace(fullName) != ""},
		Role:     string(role),
		TenantID: tenant,
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
	var tenantID *int64
	if row.TenantID.Valid {
		tenantID = &row.TenantID.Int64
	}
	fullName := ""
	if row.FullName.Valid {
		fullName = row.FullName.String
	}
	return User{
		ID:        row.ID,
		Username:  row.Username,
		Password:  row.Password,
		FullName:  fullName,
		Role:      Role(row.Role),
		TenantID:  tenantID,
		CreatedAt: row.CreatedAt.Time,
	}
}
