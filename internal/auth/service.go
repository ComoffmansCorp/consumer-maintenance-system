package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	platformauth "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/auth"
	platformcache "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/cache"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrSuperAdminExists    = errors.New("SUPER_ADMIN already exists")
	ErrUsernameInUse       = errors.New("username is already in use")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type Service struct {
	repo   *Repository
	tokens *platformauth.Service
	// cache is used to revoke an access token's jti on logout so it stops
	// working immediately instead of riding out its own TTL (see
	// revokeAccessToken). Nil is a valid value -- e.g. cmd/seed never logs
	// anyone out, so it wires this domain up without a cache client at all.
	cache *platformcache.Client
}

func NewService(repo *Repository, tokens *platformauth.Service, cache *platformcache.Client) *Service {
	return &Service{repo: repo, tokens: tokens, cache: cache}
}

func (s *Service) BootstrapSuperAdmin(ctx context.Context, req BootstrapSuperAdminRequest) (AuthResponse, error) {
	if err := validateBootstrap(req); err != nil {
		return AuthResponse{}, err
	}

	count, err := s.repo.CountByRole(ctx, RoleSuperAdmin)
	if err != nil {
		return AuthResponse{}, err
	}
	if count > 0 {
		return AuthResponse{}, ErrSuperAdminExists
	}

	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return AuthResponse{}, err
	}
	if exists {
		return AuthResponse{}, ErrUsernameInUse
	}

	passwordHash, err := platformauth.HashPassword(req.Password)
	if err != nil {
		return AuthResponse{}, err
	}

	user, err := s.repo.Create(ctx, strings.TrimSpace(req.Username), passwordHash, req.FullName, RoleSuperAdmin)
	if err != nil {
		return AuthResponse{}, err
	}

	return s.buildAuthResponse(ctx, user)
}

// RegisterClient and RegisterMaster create marketplace users. Username
// uniqueness is checked against the shared username pool (shared with
// SUPER_ADMIN), matching the ux_users_username_lower index.
func (s *Service) RegisterClient(ctx context.Context, req RegisterRequest) (AuthResponse, error) {
	return s.register(ctx, req, RoleClient)
}

func (s *Service) RegisterMaster(ctx context.Context, req RegisterRequest) (AuthResponse, error) {
	return s.register(ctx, req, RoleMaster)
}

func (s *Service) register(ctx context.Context, req RegisterRequest, role Role) (AuthResponse, error) {
	if err := validateRegistration(req); err != nil {
		return AuthResponse{}, err
	}

	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return AuthResponse{}, err
	}
	if exists {
		return AuthResponse{}, ErrUsernameInUse
	}

	passwordHash, err := platformauth.HashPassword(req.Password)
	if err != nil {
		return AuthResponse{}, err
	}

	user, err := s.repo.Create(ctx, strings.TrimSpace(req.Username), passwordHash, req.FullName, role)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("register user: %w", err)
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (AuthResponse, error) {
	if err := validateLogin(req); err != nil {
		return AuthResponse{}, err
	}

	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthResponse{}, ErrInvalidCredentials
		}
		return AuthResponse{}, err
	}

	if err := platformauth.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		return AuthResponse{}, ErrInvalidCredentials
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *Service) Refresh(ctx context.Context, req RefreshRequest) (AuthResponse, error) {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return AuthResponse{}, ErrInvalidRefreshToken
	}

	tokenHash := hashToken(req.RefreshToken)
	stored, err := s.repo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthResponse{}, ErrInvalidRefreshToken
		}
		return AuthResponse{}, err
	}
	if time.Now().UTC().After(stored.ExpiresAt) {
		return AuthResponse{}, ErrInvalidRefreshToken
	}

	user, err := s.repo.GetByID(ctx, stored.UserID)
	if err != nil {
		return AuthResponse{}, err
	}

	if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return AuthResponse{}, err
	}

	return s.buildAuthResponse(ctx, user)
}

// GetUser returns basic user info by id, trusting the caller already
// established the id is relevant. Used by other domains (via a UserStore
// port) to resolve display names.
func (s *Service) GetUser(ctx context.Context, id int64) (UserDTO, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserDTO{}, fmt.Errorf("user not found")
		}
		return UserDTO{}, err
	}
	return ToUserDTO(user), nil
}

// Logout revokes the refresh token in Postgres (as before) and, best-effort,
// the access token's jti in Redis (accessToken is whatever the client sent
// in its Authorization header, possibly empty -- /api/auth/logout is a
// public route, see router.go, so it isn't guaranteed to carry one).
func (s *Service) Logout(ctx context.Context, req LogoutRequest, accessToken string) error {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return ErrInvalidRefreshToken
	}
	if err := s.repo.RevokeRefreshToken(ctx, hashToken(req.RefreshToken)); err != nil {
		return err
	}
	s.revokeAccessToken(ctx, accessToken)
	return nil
}

// revokeAccessToken records the access token's jti in Redis with a TTL
// matching its own remaining lifetime, so middleware.JWTAuth's revocation
// check (internal/platform/middleware/auth.go) rejects it immediately
// instead of letting it ride out its own expiry after logout. Anything that
// keeps this from happening -- no cache configured, no/invalid/expired
// token, a Redis hiccup -- is swallowed: the refresh token is already
// revoked above, which is the part that actually matters for correctness;
// this is defense in depth on top of it, not the thing logout depends on.
func (s *Service) revokeAccessToken(ctx context.Context, accessToken string) {
	if s.cache == nil || strings.TrimSpace(accessToken) == "" {
		return
	}
	claims, err := s.tokens.ParseAccessToken(accessToken)
	if err != nil || claims.ID == "" || claims.ExpiresAt == nil {
		return
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return
	}
	_ = s.cache.Set(ctx, "revoked:"+claims.ID, "1", ttl)
}

func (s *Service) buildAuthResponse(ctx context.Context, user User) (AuthResponse, error) {
	accessToken, expiresIn, err := s.tokens.IssueAccessToken(user.ID, string(user.Role))
	if err != nil {
		return AuthResponse{}, err
	}

	refreshToken, err := s.issueRefreshToken(ctx, user.ID)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		UserID:       user.ID,
		FullName:     user.FullName,
		Role:         user.Role,
	}, nil
}

func (s *Service) issueRefreshToken(ctx context.Context, userID int64) (string, error) {
	raw, err := randomToken()
	if err != nil {
		return "", err
	}
	expiresAt := time.Now().UTC().Add(s.tokens.RefreshTTL())
	if err := s.repo.CreateRefreshToken(ctx, userID, hashToken(raw), expiresAt); err != nil {
		return "", err
	}
	return raw, nil
}

func validateRegistration(req RegisterRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(req.FullName) == "" {
		return fmt.Errorf("full name is required")
	}
	if len(req.Password) < 6 || len(req.Password) > 128 {
		return fmt.Errorf("password length must be between 6 and 128")
	}
	return nil
}

func validateBootstrap(req BootstrapSuperAdminRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if len(req.Password) < 6 || len(req.Password) > 128 {
		return fmt.Errorf("password length must be between 6 and 128")
	}
	return nil
}

func validateLogin(req LoginRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if len(req.Password) < 6 || len(req.Password) > 128 {
		return fmt.Errorf("password length must be between 6 and 128")
	}
	return nil
}

func randomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
