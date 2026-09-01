package marketplace

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

var (
	ErrCategoryNotFound     = errors.New("category not found")
	ErrServiceNotFound      = errors.New("service not found")
	ErrRequestNotFound      = errors.New("request not found")
	ErrNotSpecialized       = errors.New("master is not specialized in this service")
	ErrInvalidTransition    = errors.New("invalid request status transition")
	ErrNotOwner             = errors.New("only the client who created the request can perform this action")
	ErrNotAssignedMaster    = errors.New("only the assigned master can perform this action")
	ErrCancelReasonRequired = errors.New("cancel reason is required")
	ErrDescriptionRequired  = errors.New("description is required")
	ErrAddressRequired      = errors.New("address is required")
)

type Service struct {
	repo  *Repository
	users UserStore
}

func NewService(repo *Repository, users UserStore) *Service {
	return &Service{repo: repo, users: users}
}

// --- catalog: browsing is public, creation is SUPER_ADMIN-only (enforced
// in the handler via RequireRoles, not repeated here) ---

func (s *Service) ListCategories(ctx context.Context) ([]CategoryDTO, error) {
	categories, err := s.repo.ListActiveCategories(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]CategoryDTO, 0, len(categories))
	for _, c := range categories {
		out = append(out, ToCategoryDTO(c))
	}
	return out, nil
}

func (s *Service) CreateCategory(ctx context.Context, name string) (CategoryDTO, error) {
	if strings.TrimSpace(name) == "" {
		return CategoryDTO{}, fmt.Errorf("category name is required")
	}
	c, err := s.repo.CreateCategory(ctx, strings.TrimSpace(name))
	if err != nil {
		return CategoryDTO{}, err
	}
	return ToCategoryDTO(c), nil
}

func (s *Service) ListServices(ctx context.Context, categoryID int64) ([]ServiceDTO, error) {
	offerings, err := s.repo.ListActiveServices(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	out := make([]ServiceDTO, 0, len(offerings))
	for _, o := range offerings {
		out = append(out, ToServiceDTO(o))
	}
	return out, nil
}

func (s *Service) CreateService(ctx context.Context, req CreateOfferingRequest) (ServiceDTO, error) {
	if strings.TrimSpace(req.Name) == "" {
		return ServiceDTO{}, fmt.Errorf("service name is required")
	}
	if _, err := s.repo.GetCategoryByID(ctx, req.CategoryID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ServiceDTO{}, ErrCategoryNotFound
		}
		return ServiceDTO{}, err
	}
	o, err := s.repo.CreateService(ctx, req.CategoryID, strings.TrimSpace(req.Name), req.Description)
	if err != nil {
		return ServiceDTO{}, err
	}
	return ToServiceDTO(o), nil
}

// --- master profile ---

func (s *Service) GetMasterProfile(ctx context.Context, masterUserID int64) (MasterProfileDTO, error) {
	profile, err := s.repo.GetMasterProfile(ctx, masterUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return MasterProfileDTO{SpecializationIDs: []int64{}}, nil
		}
		return MasterProfileDTO{}, err
	}
	ids, err := s.repo.ListMasterSpecializationIDs(ctx, masterUserID)
	if err != nil {
		return MasterProfileDTO{}, err
	}
	return MasterProfileDTO{City: profile.City, Bio: profile.Bio, SpecializationIDs: ids}, nil
}

// UpdateMasterProfile upserts the 1:1 profile row lazily on first write --
// deliberately not created at registration time, so the foundational auth
// domain never has to depend on marketplace.
func (s *Service) UpdateMasterProfile(ctx context.Context, masterUserID int64, req UpdateMasterProfileRequest) (MasterProfileDTO, error) {
	if _, err := s.repo.UpsertMasterProfile(ctx, masterUserID, strings.TrimSpace(req.City), strings.TrimSpace(req.Bio)); err != nil {
		return MasterProfileDTO{}, err
	}

	validIDs := make([]int64, 0, len(req.SpecializationIDs))
	for _, id := range req.SpecializationIDs {
		if _, err := s.repo.GetServiceByID(ctx, id); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return MasterProfileDTO{}, ErrServiceNotFound
			}
			return MasterProfileDTO{}, err
		}
		validIDs = append(validIDs, id)
	}
	if err := s.repo.ReplaceMasterSpecializations(ctx, masterUserID, validIDs); err != nil {
		return MasterProfileDTO{}, err
	}

	return s.GetMasterProfile(ctx, masterUserID)
}

// --- requests (заявки) ---

func (s *Service) CreateRequest(ctx context.Context, clientID int64, req CreateRequestRequest) (RequestDTO, error) {
	if strings.TrimSpace(req.Description) == "" {
		return RequestDTO{}, ErrDescriptionRequired
	}
	if strings.TrimSpace(req.AddressText) == "" {
		return RequestDTO{}, ErrAddressRequired
	}
	offering, err := s.repo.GetServiceByID(ctx, req.ServiceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RequestDTO{}, ErrServiceNotFound
		}
		return RequestDTO{}, err
	}
	// Active is checked explicitly, not just relied on via ListServices
	// filtering it out of the catalog browse -- a stale serviceId from a
	// client that cached the catalog before deactivation must still be
	// rejected here, not just hidden from view.
	if !offering.Active {
		return RequestDTO{}, ErrServiceNotFound
	}

	r, err := s.repo.CreateRequest(ctx, clientID, req.ServiceID, strings.TrimSpace(req.Description), strings.TrimSpace(req.AddressText), req.Latitude, req.Longitude)
	if err != nil {
		return RequestDTO{}, fmt.Errorf("create request: %w", err)
	}
	return s.enrich(ctx, r), nil
}

func (s *Service) Get(ctx context.Context, id int64) (RequestDTO, error) {
	r, err := s.getOrNotFound(ctx, id)
	if err != nil {
		return RequestDTO{}, err
	}
	return s.enrich(ctx, r), nil
}

func (s *Service) ListMine(ctx context.Context, clientID int64, page, pageSize int) (httpx.Page[RequestDTO], error) {
	limit, offset := int32(pageSize), httpx.Offset(page, pageSize)
	items, err := s.repo.ListByClient(ctx, clientID, limit, offset)
	if err != nil {
		return httpx.Page[RequestDTO]{}, err
	}
	total, err := s.repo.CountByClient(ctx, clientID)
	if err != nil {
		return httpx.Page[RequestDTO]{}, err
	}
	return httpx.NewPage(s.enrichAll(ctx, items), page, pageSize, total), nil
}

func (s *Service) ListClaimedByMaster(ctx context.Context, masterID int64, page, pageSize int) (httpx.Page[RequestDTO], error) {
	limit, offset := int32(pageSize), httpx.Offset(page, pageSize)
	items, err := s.repo.ListByMaster(ctx, masterID, limit, offset)
	if err != nil {
		return httpx.Page[RequestDTO]{}, err
	}
	total, err := s.repo.CountByMaster(ctx, masterID)
	if err != nil {
		return httpx.Page[RequestDTO]{}, err
	}
	return httpx.NewPage(s.enrichAll(ctx, items), page, pageSize, total), nil
}

// ListOpenForMaster returns the claimable pool, already restricted server-side
// to services matching the master's own specializations -- a master never
// even sees requests for work outside what they're specialized in.
func (s *Service) ListOpenForMaster(ctx context.Context, masterUserID int64, page, pageSize int) (httpx.Page[RequestDTO], error) {
	limit, offset := int32(pageSize), httpx.Offset(page, pageSize)
	items, err := s.repo.ListOpenForMaster(ctx, masterUserID, limit, offset)
	if err != nil {
		return httpx.Page[RequestDTO]{}, err
	}
	total, err := s.repo.CountOpenForMaster(ctx, masterUserID)
	if err != nil {
		return httpx.Page[RequestDTO]{}, err
	}
	return httpx.NewPage(s.enrichAll(ctx, items), page, pageSize, total), nil
}

// Claim is a hard-gated action: specialization is checked here against the
// service the request is actually for, on the server, regardless of what
// the client filtered/showed -- an electrician-equivalent master cannot
// claim a plumbing request even by calling the API directly.
func (s *Service) Claim(ctx context.Context, id, masterUserID int64) (RequestDTO, error) {
	current, err := s.getOrNotFound(ctx, id)
	if err != nil {
		return RequestDTO{}, err
	}
	if current.Status != RequestStatusOpen {
		return RequestDTO{}, ErrInvalidTransition
	}

	has, err := s.repo.MasterHasSpecialization(ctx, masterUserID, current.OfferingID)
	if err != nil {
		return RequestDTO{}, err
	}
	if !has {
		return RequestDTO{}, ErrNotSpecialized
	}

	r, err := s.repo.Claim(ctx, id, masterUserID)
	if err != nil {
		return RequestDTO{}, fmt.Errorf("claim request: %w", err)
	}
	return s.enrich(ctx, r), nil
}

func (s *Service) Complete(ctx context.Context, id, masterUserID int64) (RequestDTO, error) {
	current, err := s.getOrNotFound(ctx, id)
	if err != nil {
		return RequestDTO{}, err
	}
	if current.MasterID == nil || *current.MasterID != masterUserID {
		return RequestDTO{}, ErrNotAssignedMaster
	}
	if current.Status != RequestStatusInProgress {
		return RequestDTO{}, ErrInvalidTransition
	}

	r, err := s.repo.Complete(ctx, id)
	if err != nil {
		return RequestDTO{}, fmt.Errorf("complete request: %w", err)
	}
	return s.enrich(ctx, r), nil
}

// Cancel is allowed for the owning client at any point before completion, or
// for the assigned master while they hold the claim -- either side can back
// out of a still-open engagement.
func (s *Service) Cancel(ctx context.Context, id, userID int64, reason string) (RequestDTO, error) {
	if strings.TrimSpace(reason) == "" {
		return RequestDTO{}, ErrCancelReasonRequired
	}
	current, err := s.getOrNotFound(ctx, id)
	if err != nil {
		return RequestDTO{}, err
	}
	isOwner := current.ClientID == userID
	isAssignedMaster := current.MasterID != nil && *current.MasterID == userID
	if !isOwner && !isAssignedMaster {
		return RequestDTO{}, ErrNotOwner
	}
	if current.Status == RequestStatusCompleted || current.Status == RequestStatusCanceled {
		return RequestDTO{}, ErrInvalidTransition
	}

	r, err := s.repo.Cancel(ctx, id, reason)
	if err != nil {
		return RequestDTO{}, fmt.Errorf("cancel request: %w", err)
	}
	return s.enrich(ctx, r), nil
}

func (s *Service) getOrNotFound(ctx context.Context, id int64) (Request, error) {
	r, err := s.repo.GetRequestByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Request{}, ErrRequestNotFound
		}
		return Request{}, err
	}
	return r, nil
}

func (s *Service) enrich(ctx context.Context, r Request) RequestDTO {
	dto := ToRequestDTO(r)
	if offering, err := s.repo.GetServiceByID(ctx, r.OfferingID); err == nil {
		dto.ServiceName = offering.Name
		if category, err := s.repo.GetCategoryByID(ctx, offering.CategoryID); err == nil {
			dto.CategoryName = category.Name
		}
	}
	if s.users != nil {
		if client, err := s.users.GetUser(ctx, r.ClientID); err == nil {
			dto.ClientName = client.FullName
		}
		if r.MasterID != nil {
			if master, err := s.users.GetUser(ctx, *r.MasterID); err == nil {
				dto.MasterName = master.FullName
			}
		}
	}
	return dto
}

func (s *Service) enrichAll(ctx context.Context, items []Request) []RequestDTO {
	out := make([]RequestDTO, 0, len(items))
	for _, r := range items {
		out = append(out, s.enrich(ctx, r))
	}
	return out
}
