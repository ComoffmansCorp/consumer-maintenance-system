package consumer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

var (
	ErrConsumerNotFound = errors.New("consumer not found")
	ErrInvalidType       = errors.New("invalid consumer type")
	ErrNameRequired      = errors.New("consumer name is required")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, tenantID int64, req CreateRequest) (DTO, error) {
	if err := validate(req.Name, req.Type); err != nil {
		return DTO{}, err
	}
	c, err := s.repo.Create(ctx, tenantID, req.Name, req.Type, req.Description)
	if err != nil {
		return DTO{}, fmt.Errorf("create consumer: %w", err)
	}
	return ToDTO(c), nil
}

func (s *Service) Get(ctx context.Context, id, tenantID int64) (DTO, error) {
	c, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO{}, ErrConsumerNotFound
		}
		return DTO{}, err
	}
	return ToDTO(c), nil
}

func (s *Service) List(ctx context.Context, tenantID int64, search string, page, pageSize int) (httpx.Page[DTO], error) {
	limit := int32(pageSize)
	offset := httpx.Offset(page, pageSize)

	consumers, err := s.repo.List(ctx, tenantID, search, limit, offset)
	if err != nil {
		return httpx.Page[DTO]{}, err
	}
	total, err := s.repo.Count(ctx, tenantID, search)
	if err != nil {
		return httpx.Page[DTO]{}, err
	}

	items := make([]DTO, 0, len(consumers))
	for _, c := range consumers {
		items = append(items, ToDTO(c))
	}
	return httpx.NewPage(items, page, pageSize, total), nil
}

func (s *Service) Update(ctx context.Context, id, tenantID int64, req UpdateRequest) (DTO, error) {
	if err := validate(req.Name, req.Type); err != nil {
		return DTO{}, err
	}
	c, err := s.repo.Update(ctx, id, tenantID, req.Name, req.Type, req.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO{}, ErrConsumerNotFound
		}
		return DTO{}, fmt.Errorf("update consumer: %w", err)
	}
	return ToDTO(c), nil
}

func (s *Service) Deactivate(ctx context.Context, id, tenantID int64) (DTO, error) {
	c, err := s.repo.Deactivate(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO{}, ErrConsumerNotFound
		}
		return DTO{}, fmt.Errorf("deactivate consumer: %w", err)
	}
	return ToDTO(c), nil
}

// Exists is used by other domains (via the ConsumerStore port) to validate a
// consumer reference without a cross-domain join.
func (s *Service) Exists(ctx context.Context, id, tenantID int64) (bool, error) {
	return s.repo.Exists(ctx, id, tenantID)
}

// Summary returns the minimal projection other domains need to enrich their
// own responses (e.g. address listing a consumer name) without a join.
func (s *Service) Summary(ctx context.Context, id, tenantID int64) (Consumer, error) {
	c, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Consumer{}, ErrConsumerNotFound
		}
		return Consumer{}, err
	}
	return c, nil
}

func validate(name string, typ Type) error {
	if strings.TrimSpace(name) == "" {
		return ErrNameRequired
	}
	if !typ.Valid() {
		return ErrInvalidType
	}
	return nil
}
