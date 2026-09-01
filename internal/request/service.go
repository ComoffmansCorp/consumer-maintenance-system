package request

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/broker"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

var (
	ErrRequestNotFound      = errors.New("request not found")
	ErrServiceNotFound      = errors.New("service not found")
	ErrNotSpecialized       = errors.New("master is not specialized in this service")
	ErrInvalidTransition    = errors.New("invalid request status transition")
	ErrNotOwner             = errors.New("only the client who created the request can perform this action")
	ErrNotAssignedMaster    = errors.New("only the assigned master can perform this action")
	ErrCancelReasonRequired = errors.New("cancel reason is required")
	ErrDescriptionRequired  = errors.New("description is required")
	ErrAddressRequired      = errors.New("address is required")
	ErrPriceInvalid         = errors.New("price must be greater than zero")
	ErrOfferNotFound        = errors.New("offer not found")
	ErrOfferNotPending      = errors.New("offer is not pending")
	ErrOfferAlreadyExists   = errors.New("master has already made an offer on this request")
)

type Service struct {
	repo     *Repository
	catalog  CatalogPort
	masters  SpecializationPort
	txRunner TxRunner
	events   *broker.Bus
}

func NewService(repo *Repository, catalog CatalogPort, masters SpecializationPort, txRunner TxRunner, events *broker.Bus) *Service {
	return &Service{repo: repo, catalog: catalog, masters: masters, txRunner: txRunner, events: events}
}

// --- requests ---

func (s *Service) CreateRequest(ctx context.Context, clientID int64, req CreateRequestRequest) (RequestDTO, error) {
	if strings.TrimSpace(req.Description) == "" {
		return RequestDTO{}, ErrDescriptionRequired
	}
	if strings.TrimSpace(req.AddressText) == "" {
		return RequestDTO{}, ErrAddressRequired
	}
	svc, found, err := s.catalog.GetService(ctx, req.ServiceID)
	if err != nil {
		return RequestDTO{}, err
	}
	// Active is checked explicitly, not just relied on via the catalog
	// browse filtering it out -- a stale serviceId from a client that
	// cached the catalog before deactivation must still be rejected here.
	if !found || !svc.Active {
		return RequestDTO{}, ErrServiceNotFound
	}

	r, err := s.repo.CreateRequest(ctx, clientID, req.ServiceID, strings.TrimSpace(req.Description), strings.TrimSpace(req.AddressText), req.Latitude, req.Longitude)
	if err != nil {
		return RequestDTO{}, fmt.Errorf("create request: %w", err)
	}
	if err := s.repo.AddStatusHistory(ctx, r.ID, "", string(StatusOpen), clientID, ""); err != nil {
		return RequestDTO{}, fmt.Errorf("record open history: %w", err)
	}
	return s.enrich(ctx, r, true), nil
}

func (s *Service) Get(ctx context.Context, id int64) (RequestDTO, error) {
	r, err := s.getOrNotFound(ctx, id)
	if err != nil {
		return RequestDTO{}, err
	}
	return s.enrich(ctx, r, true), nil
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

// ListOpenForMaster returns the biddable pool, already restricted
// server-side to services matching the master's own specializations -- a
// master never even sees requests for work outside what they're specialized
// in.
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

func (s *Service) ListAdmin(ctx context.Context, status string, page, pageSize int) (httpx.Page[RequestDTO], error) {
	limit, offset := int32(pageSize), httpx.Offset(page, pageSize)
	items, err := s.repo.ListAdmin(ctx, status, limit, offset)
	if err != nil {
		return httpx.Page[RequestDTO]{}, err
	}
	total, err := s.repo.CountAdmin(ctx, status)
	if err != nil {
		return httpx.Page[RequestDTO]{}, err
	}
	return httpx.NewPage(s.enrichAll(ctx, items), page, pageSize, total), nil
}

// --- offers ---

// SubmitOffer is a hard-gated action: specialization is checked here
// against the service the request is actually for, on the server,
// regardless of what the client filtered/showed -- an electrician-equivalent
// master cannot bid on a plumbing request even by calling the API directly.
func (s *Service) SubmitOffer(ctx context.Context, requestID, masterUserID int64, req SubmitOfferRequest) (OfferDTO, error) {
	if req.Price <= 0 {
		return OfferDTO{}, ErrPriceInvalid
	}
	current, err := s.getOrNotFound(ctx, requestID)
	if err != nil {
		return OfferDTO{}, err
	}
	if current.Status != StatusOpen {
		return OfferDTO{}, ErrInvalidTransition
	}

	has, err := s.masters.HasSpecialization(ctx, masterUserID, current.ServiceID)
	if err != nil {
		return OfferDTO{}, err
	}
	if !has {
		return OfferDTO{}, ErrNotSpecialized
	}

	o, err := s.repo.CreateOffer(ctx, requestID, masterUserID, req.Price, strings.TrimSpace(req.Comment))
	if err != nil {
		if isUniqueViolation(err) {
			return OfferDTO{}, ErrOfferAlreadyExists
		}
		return OfferDTO{}, fmt.Errorf("submit offer: %w", err)
	}
	return ToOfferDTO(o), nil
}

// ListOffers is restricted to the request owner: a master only ever sees
// their own offer, never a competitor's price, via the request DTO.
func (s *Service) ListOffers(ctx context.Context, requestID, clientID int64) ([]OfferDTO, error) {
	current, err := s.getOrNotFound(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if current.ClientID != clientID {
		return nil, ErrNotOwner
	}
	offers, err := s.repo.ListOffersByRequest(ctx, requestID)
	if err != nil {
		return nil, err
	}
	out := make([]OfferDTO, 0, len(offers))
	for _, o := range offers {
		out = append(out, ToOfferDTO(o))
	}
	return out, nil
}

// AcceptOffer is transactional: the accepted offer, every other still-PENDING
// offer on the same request, and the request itself all move together --
// either everything commits or nothing does.
func (s *Service) AcceptOffer(ctx context.Context, requestID, offerID, clientID int64) (RequestDTO, error) {
	current, err := s.getOrNotFound(ctx, requestID)
	if err != nil {
		return RequestDTO{}, err
	}
	if current.ClientID != clientID {
		return RequestDTO{}, ErrNotOwner
	}
	if current.Status != StatusOpen {
		return RequestDTO{}, ErrInvalidTransition
	}

	offer, err := s.repo.GetOfferByID(ctx, offerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RequestDTO{}, ErrOfferNotFound
		}
		return RequestDTO{}, err
	}
	if offer.RequestID != requestID {
		return RequestDTO{}, ErrOfferNotFound
	}
	if offer.Status != OfferPending {
		return RequestDTO{}, ErrOfferNotPending
	}

	var updated ServiceRequest
	err = s.txRunner.WithinTransaction(ctx, func(txCtx context.Context) error {
		if _, err := s.repo.AcceptOffer(txCtx, offerID); err != nil {
			return err
		}
		if err := s.repo.RejectOtherOffers(txCtx, requestID, offerID); err != nil {
			return err
		}
		r, err := s.repo.Assign(txCtx, requestID, offer.MasterID, offer.Price)
		if err != nil {
			return err
		}
		if err := s.repo.AddStatusHistory(txCtx, requestID, string(StatusOpen), string(StatusAssigned), clientID, ""); err != nil {
			return err
		}
		updated = r
		return nil
	})
	if err != nil {
		return RequestDTO{}, fmt.Errorf("accept offer: %w", err)
	}

	s.events.Publish(ctx, broker.Event{
		Type: broker.EventRequestAssigned,
		Payload: map[string]any{
			"request_id": updated.ID,
			"master_id":  offer.MasterID,
			"price":      offer.Price,
		},
	})

	return s.enrich(ctx, updated, true), nil
}

func (s *Service) CompleteRequest(ctx context.Context, requestID, masterUserID int64) (RequestDTO, error) {
	current, err := s.getOrNotFound(ctx, requestID)
	if err != nil {
		return RequestDTO{}, err
	}
	if current.MasterID == nil || *current.MasterID != masterUserID {
		return RequestDTO{}, ErrNotAssignedMaster
	}
	if current.Status != StatusAssigned {
		return RequestDTO{}, ErrInvalidTransition
	}

	r, err := s.repo.Complete(ctx, requestID)
	if err != nil {
		return RequestDTO{}, fmt.Errorf("complete request: %w", err)
	}
	if err := s.repo.AddStatusHistory(ctx, requestID, string(StatusAssigned), string(StatusCompleted), masterUserID, ""); err != nil {
		return RequestDTO{}, fmt.Errorf("record complete history: %w", err)
	}

	s.events.Publish(ctx, broker.Event{
		Type:    broker.EventRequestCompleted,
		Payload: map[string]any{"request_id": r.ID},
	})

	return s.enrich(ctx, r, true), nil
}

// CancelRequest is allowed for the owning client at any point before
// completion, or for the assigned master while they hold the assignment --
// either side can back out of a still-open engagement.
func (s *Service) CancelRequest(ctx context.Context, requestID, userID int64, reason string) (RequestDTO, error) {
	if strings.TrimSpace(reason) == "" {
		return RequestDTO{}, ErrCancelReasonRequired
	}
	current, err := s.getOrNotFound(ctx, requestID)
	if err != nil {
		return RequestDTO{}, err
	}
	isOwner := current.ClientID == userID
	isAssignedMaster := current.MasterID != nil && *current.MasterID == userID
	if !isOwner && !isAssignedMaster {
		return RequestDTO{}, ErrNotOwner
	}
	if current.Status == StatusCompleted || current.Status == StatusCanceled {
		return RequestDTO{}, ErrInvalidTransition
	}

	r, err := s.repo.Cancel(ctx, requestID, reason)
	if err != nil {
		return RequestDTO{}, fmt.Errorf("cancel request: %w", err)
	}
	if err := s.repo.AddStatusHistory(ctx, requestID, string(current.Status), string(StatusCanceled), userID, reason); err != nil {
		return RequestDTO{}, fmt.Errorf("record cancel history: %w", err)
	}

	s.events.Publish(ctx, broker.Event{
		Type:    broker.EventRequestCanceled,
		Payload: map[string]any{"request_id": r.ID},
	})

	return s.enrich(ctx, r, true), nil
}

// --- favorites ---

func (s *Service) AddFavorite(ctx context.Context, clientID, masterID int64) error {
	return s.repo.AddFavorite(ctx, clientID, masterID)
}

func (s *Service) RemoveFavorite(ctx context.Context, clientID, masterID int64) error {
	return s.repo.RemoveFavorite(ctx, clientID, masterID)
}

func (s *Service) ListFavorites(ctx context.Context, clientID int64) ([]FavoriteDTO, error) {
	favorites, err := s.repo.ListFavoritesByClient(ctx, clientID)
	if err != nil {
		return nil, err
	}
	out := make([]FavoriteDTO, 0, len(favorites))
	for _, f := range favorites {
		out = append(out, ToFavoriteDTO(f))
	}
	return out, nil
}

// --- ports for other domains (review, payment, chat), implemented here and
// wrapped by wiring adapters that satisfy each consumer's own port
// interface -- request never imports them back. ---

// GetParticipants returns who a request belongs to, for authorization
// checks in the payment and chat domains.
func (s *Service) GetParticipants(ctx context.Context, requestID int64) (clientID int64, masterID *int64, err error) {
	r, err := s.getOrNotFound(ctx, requestID)
	if err != nil {
		return 0, nil, err
	}
	return r.ClientID, r.MasterID, nil
}

// GetCompletedForReview validates that requestID belongs to clientID and is
// COMPLETED, returning the assigned master to attribute the review to.
func (s *Service) GetCompletedForReview(ctx context.Context, requestID, clientID int64) (masterID int64, err error) {
	r, err := s.getOrNotFound(ctx, requestID)
	if err != nil {
		return 0, err
	}
	if r.ClientID != clientID {
		return 0, ErrNotOwner
	}
	if r.Status != StatusCompleted || r.MasterID == nil {
		return 0, ErrInvalidTransition
	}
	return *r.MasterID, nil
}

func (s *Service) getOrNotFound(ctx context.Context, id int64) (ServiceRequest, error) {
	r, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ServiceRequest{}, ErrRequestNotFound
		}
		return ServiceRequest{}, err
	}
	return r, nil
}

func (s *Service) enrich(ctx context.Context, r ServiceRequest, withHistory bool) RequestDTO {
	dto := ToRequestDTO(r)
	if info, found, err := s.catalog.GetService(ctx, r.ServiceID); err == nil && found {
		dto.ServiceName = info.Name
	}
	if withHistory {
		if history, err := s.repo.ListStatusHistory(ctx, r.ID); err == nil {
			entries := make([]StatusHistoryEntryDTO, 0, len(history))
			for _, h := range history {
				entries = append(entries, ToStatusHistoryDTO(h))
			}
			dto.History = entries
		}
	}
	return dto
}

func (s *Service) enrichAll(ctx context.Context, items []ServiceRequest) []RequestDTO {
	out := make([]RequestDTO, 0, len(items))
	for _, r := range items {
		out = append(out, s.enrich(ctx, r, false))
	}
	return out
}

// isUniqueViolation detects a Postgres unique-constraint violation (SQLSTATE
// 23505) -- used to turn the request_offers (request_id, master_id) unique
// index into ErrOfferAlreadyExists instead of a raw db error.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
