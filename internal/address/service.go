package address

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

var (
	ErrAddressNotFound  = errors.New("address not found")
	ErrStreetRequired   = errors.New("street is required")
	ErrHouseRequired    = errors.New("house is required")
	ErrConsumerNotFound = errors.New("consumer not found")
)

type Service struct {
	repo      *Repository
	consumers ConsumerStore
}

func NewService(repo *Repository, consumers ConsumerStore) *Service {
	return &Service{repo: repo, consumers: consumers}
}

func (s *Service) Create(ctx context.Context, tenantID int64, req CreateRequest) (DTO, error) {
	if err := validate(req.Street, req.House); err != nil {
		return DTO{}, err
	}
	if err := s.checkConsumer(ctx, tenantID, req.ConsumerID); err != nil {
		return DTO{}, err
	}

	a, err := s.repo.Create(ctx, tenantID, req.Street, req.House, req.Building, req.Apartment, req.ConsumerID)
	if err != nil {
		return DTO{}, fmt.Errorf("create address: %w", err)
	}
	return s.enrich(ctx, tenantID, a), nil
}

func (s *Service) Get(ctx context.Context, id, tenantID int64) (DTO, error) {
	a, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO{}, ErrAddressNotFound
		}
		return DTO{}, err
	}
	return s.enrich(ctx, tenantID, a), nil
}

// Exists is exposed via a port to other domains (task, act) that reference
// an address without owning its table.
func (s *Service) Exists(ctx context.Context, id, tenantID int64) (bool, error) {
	return s.repo.Exists(ctx, id, tenantID)
}

// Label returns a short human-readable address string for other domains to
// embed in their own DTOs (task list, act header) without a join.
func (s *Service) Label(ctx context.Context, id, tenantID int64) (string, error) {
	a, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrAddressNotFound
		}
		return "", err
	}
	label := a.Street + ", " + a.House
	if a.Building != "" {
		label += " корп. " + a.Building
	}
	if a.Apartment != "" {
		label += ", кв. " + a.Apartment
	}
	return label, nil
}

func (s *Service) List(ctx context.Context, tenantID int64, search string, consumerID int64, page, pageSize int) (httpx.Page[DTO], error) {
	limit := int32(pageSize)
	offset := httpx.Offset(page, pageSize)

	addresses, err := s.repo.List(ctx, tenantID, search, consumerID, limit, offset)
	if err != nil {
		return httpx.Page[DTO]{}, err
	}
	total, err := s.repo.Count(ctx, tenantID, search, consumerID)
	if err != nil {
		return httpx.Page[DTO]{}, err
	}

	items := make([]DTO, 0, len(addresses))
	for _, a := range addresses {
		items = append(items, s.enrich(ctx, tenantID, a))
	}
	return httpx.NewPage(items, page, pageSize, total), nil
}

func (s *Service) Update(ctx context.Context, id, tenantID int64, req UpdateRequest) (DTO, error) {
	if err := validate(req.Street, req.House); err != nil {
		return DTO{}, err
	}
	if err := s.checkConsumer(ctx, tenantID, req.ConsumerID); err != nil {
		return DTO{}, err
	}

	a, err := s.repo.Update(ctx, id, tenantID, req.Street, req.House, req.Building, req.Apartment, req.ConsumerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO{}, ErrAddressNotFound
		}
		return DTO{}, fmt.Errorf("update address: %w", err)
	}
	return s.enrich(ctx, tenantID, a), nil
}

func (s *Service) checkConsumer(ctx context.Context, tenantID int64, consumerID *int64) error {
	if consumerID == nil {
		return nil
	}
	exists, err := s.consumers.Exists(ctx, *consumerID, tenantID)
	if err != nil {
		return fmt.Errorf("check consumer: %w", err)
	}
	if !exists {
		return ErrConsumerNotFound
	}
	return nil
}

func (s *Service) enrich(ctx context.Context, tenantID int64, a Address) DTO {
	dto := ToDTO(a)
	if a.ConsumerID == nil {
		return dto
	}
	summary, err := s.consumers.GetSummary(ctx, *a.ConsumerID, tenantID)
	if err == nil {
		dto.ConsumerName = summary.Name
	}
	return dto
}

func validate(street, house string) error {
	if strings.TrimSpace(street) == "" {
		return ErrStreetRequired
	}
	if strings.TrimSpace(house) == "" {
		return ErrHouseRequired
	}
	return nil
}
