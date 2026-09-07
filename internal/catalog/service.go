package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	platformcache "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/cache"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrServiceNotFound  = errors.New("service not found")
	ErrNameRequired     = errors.New("name is required")
)

// catalogCacheTTL is deliberately short: this is a cache-aside convenience
// for a public, read-heavy, rarely-changing catalog, not a source of truth.
// Writes also explicitly invalidate the keys they affect (see
// invalidateCategories/invalidateServices) so an admin edit shows up
// immediately in the demo instead of waiting out the TTL.
const catalogCacheTTL = 60 * time.Second

const cacheKeyCategories = "catalog:categories"

type Service struct {
	repo   *Repository
	cache  *platformcache.Client
	logger *slog.Logger
}

// NewService takes an optional cache client -- nil disables the cache-aside
// behavior entirely (see cmd/seed, which has no need for it) and every read
// simply falls through to the repository, same as a Redis outage would.
// logger may also be nil (falls back to slog.Default()); a cache miss due to
// a Redis error is logged, never surfaced as a request failure.
func NewService(repo *Repository, cache *platformcache.Client, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{repo: repo, cache: cache, logger: logger}
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
	s.invalidateCategories(ctx)
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
	s.invalidateCategories(ctx)
	return ToCategoryDTO(c), nil
}

// ListCategories returns the active catalog as a two-level tree: every
// top-level category with its active subcategories nested inside.
// One level of nesting only, matching the schema. Cache-aside on
// cacheKeyCategories: a hit skips Postgres entirely, a miss (including any
// Redis error, treated the same as a miss) falls through and repopulates it.
func (s *Service) ListCategories(ctx context.Context) ([]CategoryDTO, error) {
	if cached, ok := s.getCached(ctx, cacheKeyCategories); ok {
		var out []CategoryDTO
		if err := json.Unmarshal([]byte(cached), &out); err == nil {
			return out, nil
		}
		s.logger.Warn("catalog cache decode failed, falling back to postgres", "key", cacheKeyCategories)
	}

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
	s.setCached(ctx, cacheKeyCategories, out)
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
	svc, err := s.repo.CreateService(ctx, req.CategoryID, strings.TrimSpace(req.Name), req.Description, req.PriceFrom, req.PriceTo, req.Unit, req.ImageURL)
	if err != nil {
		return ServiceDTO{}, fmt.Errorf("create service: %w", err)
	}
	s.invalidateServices(ctx, svc.CategoryID)
	return ToServiceDTO(svc), nil
}

func (s *Service) UpdateService(ctx context.Context, id int64, req UpdateServiceRequest) (ServiceDTO, error) {
	if strings.TrimSpace(req.Name) == "" {
		return ServiceDTO{}, ErrNameRequired
	}
	if _, err := s.getServiceOrNotFound(ctx, id); err != nil {
		return ServiceDTO{}, err
	}
	svc, err := s.repo.UpdateService(ctx, id, strings.TrimSpace(req.Name), req.Description, req.PriceFrom, req.PriceTo, req.Unit, req.Active, req.ImageURL)
	if err != nil {
		return ServiceDTO{}, fmt.Errorf("update service: %w", err)
	}
	s.invalidateServices(ctx, svc.CategoryID)
	return ToServiceDTO(svc), nil
}

// ListServices is cache-aside on a per-category key (categoryID == 0 means
// "no filter", see ListActiveServices), same rationale as ListCategories.
func (s *Service) ListServices(ctx context.Context, categoryID int64) ([]ServiceDTO, error) {
	key := servicesCacheKey(categoryID)
	if cached, ok := s.getCached(ctx, key); ok {
		var out []ServiceDTO
		if err := json.Unmarshal([]byte(cached), &out); err == nil {
			return out, nil
		}
		s.logger.Warn("catalog cache decode failed, falling back to postgres", "key", key)
	}

	rows, err := s.repo.ListActiveServices(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	out := make([]ServiceDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, ToServiceDTO(row))
	}
	s.setCached(ctx, key, out)
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

// servicesCacheKey mirrors ListActiveServices' own "0 = no filter" contract.
func servicesCacheKey(categoryID int64) string {
	if categoryID == 0 {
		return "catalog:services:all"
	}
	return "catalog:services:" + strconv.FormatInt(categoryID, 10)
}

// getCached fetches a raw cache value. false covers every reason to fall
// through to Postgres: no cache configured, a genuine miss, or a Redis
// error -- the last of which is logged, since it's the only one worth
// knowing about.
func (s *Service) getCached(ctx context.Context, key string) (string, bool) {
	if s.cache == nil {
		return "", false
	}
	value, found, err := s.cache.Get(ctx, key)
	if err != nil {
		s.logger.Warn("catalog cache unavailable, falling back to postgres", "key", key, "error", err)
		return "", false
	}
	return value, found
}

// setCached best-effort populates the cache after a Postgres read. Failures
// are logged, never propagated -- a cache write failing must not fail the
// read it's trying to speed up next time.
func (s *Service) setCached(ctx context.Context, key string, value any) {
	if s.cache == nil {
		return
	}
	raw, err := json.Marshal(value)
	if err != nil {
		s.logger.Warn("catalog cache encode failed", "key", key, "error", err)
		return
	}
	if err := s.cache.Set(ctx, key, string(raw), catalogCacheTTL); err != nil {
		s.logger.Warn("catalog cache unavailable, skipping cache write", "key", key, "error", err)
	}
}

func (s *Service) invalidateCategories(ctx context.Context) {
	if s.cache == nil {
		return
	}
	if err := s.cache.Del(ctx, cacheKeyCategories); err != nil {
		s.logger.Warn("catalog cache invalidate failed", "key", cacheKeyCategories, "error", err)
	}
}

// invalidateServices drops both the affected category's key and the "all
// services" key, since the latter's cached snapshot is now stale too.
func (s *Service) invalidateServices(ctx context.Context, categoryID int64) {
	if s.cache == nil {
		return
	}
	if err := s.cache.Del(ctx, servicesCacheKey(categoryID), servicesCacheKey(0)); err != nil {
		s.logger.Warn("catalog cache invalidate failed", "categoryId", categoryID, "error", err)
	}
}
