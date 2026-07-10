package meter

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMeterNotFound        = errors.New("meter not found")
	ErrInvalidType          = errors.New("invalid meter type")
	ErrInvalidSealState     = errors.New("invalid seal state")
	ErrSerialNumberRequired = errors.New("serial number is required")
)

type Service struct {
	repo *Repository
	acts ActStore
}

func NewService(repo *Repository, acts ActStore) *Service {
	return &Service{repo: repo, acts: acts}
}

func (s *Service) Create(ctx context.Context, actID, tenantID int64, req CreateRequest) (DTO, error) {
	if err := s.acts.EnsureInspectionAct(ctx, actID, tenantID); err != nil {
		return DTO{}, err
	}
	if err := validate(req); err != nil {
		return DTO{}, err
	}
	m, err := s.repo.Create(ctx, actID, req)
	if err != nil {
		return DTO{}, fmt.Errorf("create meter: %w", err)
	}
	return ToDTO(m), nil
}

func (s *Service) ListByAct(ctx context.Context, actID, tenantID int64) ([]DTO, error) {
	if err := s.acts.EnsureInspectionAct(ctx, actID, tenantID); err != nil {
		return nil, err
	}
	meters, err := s.repo.ListByAct(ctx, actID)
	if err != nil {
		return nil, err
	}
	out := make([]DTO, 0, len(meters))
	for _, m := range meters {
		out = append(out, ToDTO(m))
	}
	return out, nil
}

func (s *Service) Update(ctx context.Context, id, actID, tenantID int64, req UpdateRequest) (DTO, error) {
	if err := s.acts.EnsureInspectionAct(ctx, actID, tenantID); err != nil {
		return DTO{}, err
	}
	if err := validate(req); err != nil {
		return DTO{}, err
	}
	m, err := s.repo.Update(ctx, id, actID, req)
	if err != nil {
		return DTO{}, fmt.Errorf("update meter: %w", err)
	}
	return ToDTO(m), nil
}

func (s *Service) Delete(ctx context.Context, id, actID, tenantID int64) error {
	if err := s.acts.EnsureInspectionAct(ctx, actID, tenantID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id, actID)
}

func validate(req CreateRequest) error {
	if strings.TrimSpace(req.SerialNumber) == "" {
		return ErrSerialNumberRequired
	}
	if !req.Type.Valid() {
		return ErrInvalidType
	}
	if !req.SealState.Valid() {
		return ErrInvalidSealState
	}
	return nil
}
