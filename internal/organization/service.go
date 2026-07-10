package organization

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

var (
	ErrTenantNotFound   = errors.New("tenant not found")
	ErrTenantNameInUse  = errors.New("tenant name is already in use")
	ErrTenantCodeInUse  = errors.New("tenant code is already in use")
	ErrInvalidTenantPlan = errors.New("invalid tenant plan")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateTenantRequest) (TenantDTO, error) {
	if err := validateCreate(req); err != nil {
		return TenantDTO{}, err
	}

	nameExists, err := s.repo.ExistsByName(ctx, req.Name)
	if err != nil {
		return TenantDTO{}, err
	}
	if nameExists {
		return TenantDTO{}, ErrTenantNameInUse
	}

	codeExists, err := s.repo.ExistsByCode(ctx, req.Code)
	if err != nil {
		return TenantDTO{}, err
	}
	if codeExists {
		return TenantDTO{}, ErrTenantCodeInUse
	}

	tenant, err := s.repo.Create(ctx, req.Name, req.Code, req.Plan)
	if err != nil {
		return TenantDTO{}, err
	}
	return ToDTO(tenant), nil
}

func (s *Service) List(ctx context.Context) ([]TenantDTO, error) {
	tenants, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]TenantDTO, 0, len(tenants))
	for _, tenant := range tenants {
		out = append(out, ToDTO(tenant))
	}
	return out, nil
}

func (s *Service) UpdatePlan(ctx context.Context, tenantID int64, req UpdateTenantPlanRequest) (TenantDTO, error) {
	if !isValidPlan(req.Plan) {
		return TenantDTO{}, ErrInvalidTenantPlan
	}
	tenant, err := s.repo.UpdatePlan(ctx, tenantID, req.Plan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TenantDTO{}, ErrTenantNotFound
		}
		return TenantDTO{}, fmt.Errorf("update tenant plan: %w", err)
	}
	return ToDTO(tenant), nil
}

func (s *Service) Deactivate(ctx context.Context, tenantID int64) (TenantDTO, error) {
	tenant, err := s.repo.Deactivate(ctx, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TenantDTO{}, ErrTenantNotFound
		}
		return TenantDTO{}, fmt.Errorf("deactivate tenant: %w", err)
	}
	return ToDTO(tenant), nil
}

func (s *Service) GetByCode(ctx context.Context, code string) (Tenant, error) {
	tenant, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Tenant{}, ErrTenantNotFound
		}
		return Tenant{}, err
	}
	return tenant, nil
}

func validateCreate(req CreateTenantRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("tenant name is required")
	}
	if strings.TrimSpace(req.Code) == "" {
		return fmt.Errorf("tenant code is required")
	}
	if !isValidPlan(req.Plan) {
		return ErrInvalidTenantPlan
	}
	return nil
}

func isValidPlan(plan TenantPlan) bool {
	switch plan {
	case PlanFree, PlanBusiness, PlanEnterprise:
		return true
	default:
		return false
	}
}
