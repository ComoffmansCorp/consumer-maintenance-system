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
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSuperAdminExists   = errors.New("SUPER_ADMIN already exists")
	ErrUsernameInUse      = errors.New("username is already in use")
	ErrTenantNameInUse    = errors.New("tenant name is already in use")
	ErrTenantCodeInUse    = errors.New("tenant code is already in use")
	ErrUserLimitExceeded  = errors.New("tenant user limit exceeded for current plan")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type Service struct {
	repo        *Repository
	tenants     TenantStore
	tokens      *platformauth.Service
	txRunner    TxRunner
}

func NewService(repo *Repository, tenants TenantStore, tokens *platformauth.Service, txRunner TxRunner) *Service {
	return &Service{
		repo:     repo,
		tenants:  tenants,
		tokens:   tokens,
		txRunner: txRunner,
	}
}

func (s *Service) RegisterCompany(ctx context.Context, req CompanyRegistrationRequest) (AuthResponse, error) {
	if err := validateRegistration(req); err != nil {
		return AuthResponse{}, err
	}

	plan := defaultPlan(req.Plan)
	tenantCode := normalizeTenantCode(req.TenantCode)

	nameExists, err := s.tenants.ExistsByName(ctx, req.TenantName)
	if err != nil {
		return AuthResponse{}, err
	}
	if nameExists {
		return AuthResponse{}, ErrTenantNameInUse
	}

	codeExists, err := s.tenants.ExistsByCode(ctx, tenantCode)
	if err != nil {
		return AuthResponse{}, err
	}
	if codeExists {
		return AuthResponse{}, ErrTenantCodeInUse
	}

	passwordHash, err := platformauth.HashPassword(req.Password)
	if err != nil {
		return AuthResponse{}, err
	}

	var tenant TenantInfo
	var user User

	err = s.txRunner.WithinTransaction(ctx, func(txCtx context.Context) error {
		createdTenant, err := s.tenants.Create(txCtx, req.TenantName, tenantCode, plan)
		if err != nil {
			return err
		}
		tenant = createdTenant

		createdUser, err := s.repo.Create(txCtx, strings.TrimSpace(req.Username), passwordHash, req.FullName, RoleTenantAdmin, &tenant.ID)
		if err != nil {
			return err
		}
		user = createdUser
		return nil
	})
	if err != nil {
		return AuthResponse{}, err
	}

	return s.buildAuthResponse(ctx, user, &tenant)
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

	exists, err := s.repo.ExistsPlatformUser(ctx, req.Username)
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

	user, err := s.repo.Create(ctx, strings.TrimSpace(req.Username), passwordHash, req.FullName, RoleSuperAdmin, nil)
	if err != nil {
		return AuthResponse{}, err
	}

	return s.buildAuthResponse(ctx, user, nil)
}

// RegisterClient and RegisterMaster create platform-level marketplace users
// (TenantID == nil, same pattern as SUPER_ADMIN) — unlike RegisterCompany,
// no tenant is created and there's no "only once" limit; anyone can sign up.
// Username uniqueness is checked against the platform-user pool (shared with
// SUPER_ADMIN), matching the existing ux_platform_users_username_lower index.
func (s *Service) RegisterClient(ctx context.Context, req RegisterMarketplaceRequest) (AuthResponse, error) {
	return s.registerMarketplaceUser(ctx, req, RoleClient)
}

func (s *Service) RegisterMaster(ctx context.Context, req RegisterMarketplaceRequest) (AuthResponse, error) {
	return s.registerMarketplaceUser(ctx, req, RoleMaster)
}

func (s *Service) registerMarketplaceUser(ctx context.Context, req RegisterMarketplaceRequest, role Role) (AuthResponse, error) {
	if err := validateMarketplaceRegistration(req); err != nil {
		return AuthResponse{}, err
	}

	exists, err := s.repo.ExistsPlatformUser(ctx, req.Username)
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

	user, err := s.repo.Create(ctx, strings.TrimSpace(req.Username), passwordHash, req.FullName, role, nil)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("register marketplace user: %w", err)
	}

	return s.buildAuthResponse(ctx, user, nil)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (AuthResponse, error) {
	if err := validateLogin(req); err != nil {
		return AuthResponse{}, err
	}

	var user User
	var err error

	tenantCode := normalizeOptionalTenantCode(req.TenantCode)
	if tenantCode == "" {
		user, err = s.repo.GetPlatformUser(ctx, req.Username)
	} else {
		user, err = s.repo.GetTenantUser(ctx, tenantCode, req.Username)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthResponse{}, ErrInvalidCredentials
		}
		return AuthResponse{}, err
	}

	if err := platformauth.VerifyPassword(user.Password, req.Password); err != nil {
		return AuthResponse{}, ErrInvalidCredentials
	}

	var tenant *TenantInfo
	if user.TenantID != nil {
		info, err := s.tenants.GetByID(ctx, *user.TenantID)
		if err != nil {
			return AuthResponse{}, err
		}
		if !info.Active {
			return AuthResponse{}, ErrInvalidCredentials
		}
		tenant = &info
	}

	return s.buildAuthResponse(ctx, user, tenant)
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

	var tenant *TenantInfo
	if user.TenantID != nil {
		info, err := s.tenants.GetByID(ctx, *user.TenantID)
		if err != nil {
			return AuthResponse{}, err
		}
		tenant = &info
	}

	if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return AuthResponse{}, err
	}

	return s.buildAuthResponse(ctx, user, tenant)
}

func (s *Service) CreateTenantUser(ctx context.Context, tenantID int64, req CreateTenantUserRequest) (UserDTO, error) {
	if req.Role != RoleDispatcher && req.Role != RoleElectrician {
		return UserDTO{}, fmt.Errorf("role must be DISPATCHER or ELECTRICIAN")
	}
	if strings.TrimSpace(req.Username) == "" {
		return UserDTO{}, fmt.Errorf("username is required")
	}
	if len(req.Password) < 6 || len(req.Password) > 128 {
		return UserDTO{}, fmt.Errorf("password length must be between 6 and 128")
	}

	tenant, err := s.tenants.GetByID(ctx, tenantID)
	if err != nil {
		return UserDTO{}, err
	}
	count, err := s.repo.CountByTenantID(ctx, tenantID)
	if err != nil {
		return UserDTO{}, err
	}
	if int(count) >= tenant.UserLimit {
		return UserDTO{}, ErrUserLimitExceeded
	}

	exists, err := s.repo.ExistsInTenant(ctx, tenantID, req.Username)
	if err != nil {
		return UserDTO{}, err
	}
	if exists {
		return UserDTO{}, ErrUsernameInUse
	}

	passwordHash, err := platformauth.HashPassword(req.Password)
	if err != nil {
		return UserDTO{}, err
	}

	user, err := s.repo.Create(ctx, strings.TrimSpace(req.Username), passwordHash, req.FullName, req.Role, &tenantID)
	if err != nil {
		return UserDTO{}, fmt.Errorf("create tenant user: %w", err)
	}
	return ToUserDTO(user), nil
}

func (s *Service) ListTenantUsers(ctx context.Context, tenantID int64, role Role, page, pageSize int) ([]UserDTO, error) {
	limit := int32(pageSize)
	offset := int32((page - 1) * pageSize)

	users, err := s.repo.ListByTenant(ctx, tenantID, role, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]UserDTO, 0, len(users))
	for _, u := range users {
		out = append(out, ToUserDTO(u))
	}
	return out, nil
}

// GetTenantElectrician validates that the given user id belongs to the
// tenant and holds the ELECTRICIAN role. Used by the task domain (via the
// UserStore port) to validate task assignment.
func (s *Service) GetTenantElectrician(ctx context.Context, id, tenantID int64) (UserDTO, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserDTO{}, fmt.Errorf("user not found")
		}
		return UserDTO{}, err
	}
	if user.TenantID == nil || *user.TenantID != tenantID || user.Role != RoleElectrician {
		return UserDTO{}, fmt.Errorf("user is not an electrician in this tenant")
	}
	return ToUserDTO(user), nil
}

// ListTenantAdminIDs satisfies the notification domain's AdminLister port,
// used to notify every tenant admin when a task reaches a terminal state.
func (s *Service) ListTenantAdminIDs(ctx context.Context, tenantID int64) ([]int64, error) {
	admins, err := s.repo.ListByTenant(ctx, tenantID, RoleTenantAdmin, 100, 0)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(admins))
	for _, a := range admins {
		ids = append(ids, a.ID)
	}
	return ids, nil
}

// GetUser returns basic user info by id, trusting the caller already
// established the id is relevant (e.g. a task assignee). Used by the act
// domain (via the UserStore port) to print the inspector's name.
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

func (s *Service) Logout(ctx context.Context, req LogoutRequest) error {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return ErrInvalidRefreshToken
	}
	return s.repo.RevokeRefreshToken(ctx, hashToken(req.RefreshToken))
}

func (s *Service) buildAuthResponse(ctx context.Context, user User, tenant *TenantInfo) (AuthResponse, error) {
	var tenantID *int64
	var tenantCode *string
	var tenantName *string
	var tenantPlan *string

	if tenant != nil {
		tenantID = &tenant.ID
		tenantCode = &tenant.Code
		tenantName = &tenant.Name
		tenantPlan = &tenant.Plan
	} else if user.TenantID != nil {
		tenantID = user.TenantID
	}

	accessToken, expiresIn, err := s.tokens.IssueAccessToken(user.ID, string(user.Role), tenantID, tenantCode)
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
		TenantID:     tenantID,
		TenantCode:   tenantCode,
		TenantName:   tenantName,
		TenantPlan:   tenantPlan,
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

func validateRegistration(req CompanyRegistrationRequest) error {
	if strings.TrimSpace(req.TenantName) == "" {
		return fmt.Errorf("tenant name is required")
	}
	if strings.TrimSpace(req.TenantCode) == "" {
		return fmt.Errorf("tenant code is required")
	}
	if strings.TrimSpace(req.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if len(req.Password) < 6 || len(req.Password) > 128 {
		return fmt.Errorf("password length must be between 6 and 128")
	}
	return nil
}

func validateMarketplaceRegistration(req RegisterMarketplaceRequest) error {
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

func defaultPlan(plan string) string {
	plan = strings.ToUpper(strings.TrimSpace(plan))
	switch plan {
	case "BUSINESS", "ENTERPRISE":
		return plan
	default:
		return "FREE"
	}
}

func normalizeTenantCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

func normalizeOptionalTenantCode(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return ""
	}
	return normalizeTenantCode(code)
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
