package master

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

var (
	ErrProfileNotFound       = errors.New("master profile not found")
	ErrInvalidSpecialization = errors.New("service is not a valid active specialization")
)

type Service struct {
	repo    *Repository
	catalog CatalogPort
}

func NewService(repo *Repository, catalog CatalogPort) *Service {
	return &Service{repo: repo, catalog: catalog}
}

func (s *Service) GetProfile(ctx context.Context, masterUserID int64) (ProfileDTO, error) {
	profile, err := s.repo.GetProfile(ctx, masterUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ProfileDTO{SpecializationIDs: []int64{}}, nil
		}
		return ProfileDTO{}, err
	}
	ids, err := s.repo.ListSpecializationIDs(ctx, masterUserID)
	if err != nil {
		return ProfileDTO{}, err
	}
	return ToProfileDTO(profile, ids), nil
}

// UpdateProfile upserts the 1:1 profile row lazily on first write --
// deliberately not created at registration time, so the foundational auth
// domain never has to depend on master.
func (s *Service) UpdateProfile(ctx context.Context, masterUserID int64, req UpdateProfileRequest) (ProfileDTO, error) {
	if _, err := s.repo.UpsertProfile(ctx, masterUserID, strings.TrimSpace(req.City), strings.TrimSpace(req.Bio), req.AvatarURL); err != nil {
		return ProfileDTO{}, fmt.Errorf("upsert profile: %w", err)
	}

	validIDs := make([]int64, 0, len(req.SpecializationIDs))
	for _, id := range req.SpecializationIDs {
		active, err := s.catalog.IsServiceActive(ctx, id)
		if err != nil {
			return ProfileDTO{}, err
		}
		if !active {
			return ProfileDTO{}, ErrInvalidSpecialization
		}
		validIDs = append(validIDs, id)
	}
	if err := s.repo.ReplaceSpecializations(ctx, masterUserID, validIDs); err != nil {
		return ProfileDTO{}, err
	}

	return s.GetProfile(ctx, masterUserID)
}

func (s *Service) ListProfiles(ctx context.Context, page, pageSize int) (httpx.Page[ProfileDTO], error) {
	limit, offset := int32(pageSize), httpx.Offset(page, pageSize)
	items, err := s.repo.ListProfiles(ctx, limit, offset)
	if err != nil {
		return httpx.Page[ProfileDTO]{}, err
	}
	total, err := s.repo.CountProfiles(ctx)
	if err != nil {
		return httpx.Page[ProfileDTO]{}, err
	}
	out := make([]ProfileDTO, 0, len(items))
	for _, p := range items {
		ids, err := s.repo.ListSpecializationIDs(ctx, p.UserID)
		if err != nil {
			return httpx.Page[ProfileDTO]{}, err
		}
		out = append(out, ToProfileDTO(p, ids))
	}
	return httpx.NewPage(out, page, pageSize, total), nil
}

// HasSpecialization backs the SpecializationPort used by the request domain
// to hard-gate offer submission server-side.
func (s *Service) HasSpecialization(ctx context.Context, masterUserID, serviceID int64) (bool, error) {
	return s.repo.HasSpecialization(ctx, masterUserID, serviceID)
}

// RecordReview folds a new rating into the master's running average. Used by
// the review domain (via a port) right after a review is created.
func (s *Service) RecordReview(ctx context.Context, masterUserID int64, rating int) error {
	if _, err := s.repo.RecordReview(ctx, masterUserID, rating); err != nil {
		return fmt.Errorf("record review: %w", err)
	}
	return nil
}

// GetAvatarURL backs the request domain's SpecializationPort (see
// internal/request/ports.go), which also enriches an offer DTO with the
// bidding master's avatar. A master with no profile yet, or no avatar set,
// both resolve to nil -- neither is an error worth surfacing.
func (s *Service) GetAvatarURL(ctx context.Context, masterUserID int64) (*string, error) {
	profile, err := s.repo.GetProfile(ctx, masterUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return profile.AvatarURL, nil
}
