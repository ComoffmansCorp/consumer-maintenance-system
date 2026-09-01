package catalog

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrServiceNotFound  = errors.New("service not found")
	ErrNameRequired     = errors.New("name is required")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// --- categories: browsing is public, mutation is SUPER_ADMIN-only (enforced
// in the handler via RequireRoles, not repeated here) ---

func (s *Service) CreateCategory(ctx context.Context, req CreateCategoryRequest) (CategoryDTO, error) {
	if strings.TrimSpace(req.Name) == "" {
		return CategoryDTO{}, ErrNameRequired
	}
	if req.ParentCategoryID != nil {
		if _, err := s.repo.GetCategoryByID(ctx, *req.ParentCategoryID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return CategoryDTO{}, ErrCategoryNotFound
			}
			return CategoryDTO{}, err
		}
	}
	c, err := s.repo.CreateCategory(ctx, req.ParentCategoryID, strings.TrimSpace(req.Name))
	if err != nil {
		return CategoryDTO{}, fmt.Errorf("create category: %w", err)
	}
	return ToCategoryDTO(c), nil
}

func (s *Service) UpdateCategory(ctx context.Context, id int64, req UpdateCategoryRequest) (CategoryDTO, error) {
	if strings.TrimSpace(req.Name) == "" {
		return CategoryDTO{}, ErrNameRequired
	}
	if _, err := s.repo.GetCategoryByID(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CategoryDTO{}, ErrCategoryNotFound
		}
		return CategoryDTO{}, err
	}
	c, err := s.repo.UpdateCategory(ctx, id, strings.TrimSpace(req.Name), req.Active)
	if err != nil {
		return CategoryDTO{}, fmt.Errorf("update category: %w", err)
	}
	return ToCategoryDTO(c), nil
}

// ListCategories returns the active catalog as a two-level tree: every
// top-level category with its active subcategories nested inside.
// One level of nesting only, matching the schema.
func (s *Service) ListCategories(ctx context.Context) ([]CategoryDTO, error) {
	categories, err := s.repo.ListActiveCategories(ctx)
	if err != nil {
		return nil, err
	}

	byID := make(map[int64]*CategoryDTO, len(categories))
	roots := make([]*CategoryDTO, 0, len(categories))
	for _, c := range categories {
		dto := ToCategoryDTO(c)
		byID[c.ID] = &dto
	}
	for _, c := range categories {
		dto := byID[c.ID]
		if c.ParentCategoryID == nil {
			roots = append(roots, dto)
			continue
		}
		parent, ok := byID[*c.ParentCategoryID]
		if !ok {
			// Parent inactive/missing -- surface as a top-level entry rather
			// than silently dropping it from the catalog.
			roots = append(roots, dto)
			continue
		}
		parent.Subcategories = append(parent.Subcategories, *dto)
	}

	out := make([]CategoryDTO, 0, len(roots))
	for _, r := range roots {
		out = append(out, *r)
	}
	return out, nil
}

// --- services ---

func (s *Service) CreateService(ctx context.Context, req CreateServiceRequest) (ServiceDTO, error) {
	if strings.TrimSpace(req.Name) == "" {
		return ServiceDTO{}, ErrNameRequired
	}
	if _, err := s.repo.GetCategoryByID(ctx, req.CategoryID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ServiceDTO{}, ErrCategoryNotFound
		}
		return ServiceDTO{}, err
	}
	svc, err := s.repo.CreateService(ctx, req.CategoryID, strings.TrimSpace(req.Name), req.Description, req.PriceFrom, req.PriceTo, req.Unit)
	if err != nil {
		return ServiceDTO{}, fmt.Errorf("create service: %w", err)
	}
	return ToServiceDTO(svc), nil
}

func (s *Service) UpdateService(ctx context.Context, id int64, req UpdateServiceRequest) (ServiceDTO, error) {
	if strings.TrimSpace(req.Name) == "" {
		return ServiceDTO{}, ErrNameRequired
	}
	if _, err := s.getServiceOrNotFound(ctx, id); err != nil {
		return ServiceDTO{}, err
	}
	svc, err := s.repo.UpdateService(ctx, id, strings.TrimSpace(req.Name), req.Description, req.PriceFrom, req.PriceTo, req.Unit, req.Active)
	if err != nil {
		return ServiceDTO{}, fmt.Errorf("update service: %w", err)
	}
	return ToServiceDTO(svc), nil
}

func (s *Service) ListServices(ctx context.Context, categoryID int64) ([]ServiceDTO, error) {
	rows, err := s.repo.ListActiveServices(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	out := make([]ServiceDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, ToServiceDTO(row))
	}
	return out, nil
}

// GetServiceByID / IsServiceActive back the CatalogPort used by the master
// and request domains to validate serviceIds they didn't originate.
func (s *Service) GetServiceByID(ctx context.Context, id int64) (ServiceDTO, error) {
	svc, err := s.getServiceOrNotFound(ctx, id)
	if err != nil {
		return ServiceDTO{}, err
	}
	return ToServiceDTO(svc), nil
}

func (s *Service) IsServiceActive(ctx context.Context, id int64) (bool, error) {
	svc, err := s.getServiceOrNotFound(ctx, id)
	if err != nil {
		if errors.Is(err, ErrServiceNotFound) {
			return false, nil
		}
		return false, err
	}
	return svc.Active, nil
}

func (s *Service) getServiceOrNotFound(ctx context.Context, id int64) (Offering, error) {
	svc, err := s.repo.GetServiceByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Offering{}, ErrServiceNotFound
		}
		return Offering{}, err
	}
	return svc, nil
}
