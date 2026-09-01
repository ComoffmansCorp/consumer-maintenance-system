package payment

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/broker"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

const platformFeeRate = 0.10

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrRequestNotFound = errors.New("request not found")
	ErrForbidden       = errors.New("only a participant of the request can view its payment")
)

type Service struct {
	repo     *Repository
	requests RequestPort
}

func NewService(repo *Repository, requests RequestPort) *Service {
	return &Service{repo: repo, requests: requests}
}

// RegisterHandlers subscribes to the request lifecycle events that drive
// every payment state transition -- there is no direct write path from a
// handler, matching the escrow-ledger nature of the domain.
func (s *Service) RegisterHandlers(bus *broker.Bus) {
	bus.Subscribe(broker.EventRequestAssigned, s.onRequestAssigned)
	bus.Subscribe(broker.EventRequestCompleted, s.onRequestCompleted)
	bus.Subscribe(broker.EventRequestCanceled, s.onRequestCanceled)
}

func (s *Service) onRequestAssigned(ctx context.Context, event broker.Event) error {
	requestID, ok := payloadInt64(event.Payload, "request_id")
	if !ok {
		return fmt.Errorf("request.assigned event missing request_id")
	}
	price, ok := payloadFloat64(event.Payload, "price")
	if !ok {
		return fmt.Errorf("request.assigned event missing price")
	}
	fee := price * platformFeeRate
	if _, err := s.repo.Create(ctx, requestID, price, fee); err != nil {
		return fmt.Errorf("create payment for request %d: %w", requestID, err)
	}
	return nil
}

func (s *Service) onRequestCompleted(ctx context.Context, event broker.Event) error {
	requestID, ok := payloadInt64(event.Payload, "request_id")
	if !ok {
		return fmt.Errorf("request.completed event missing request_id")
	}
	if _, err := s.repo.Release(ctx, requestID); err != nil {
		return fmt.Errorf("release payment for request %d: %w", requestID, err)
	}
	return nil
}

// onRequestCanceled refunds only if a HELD payment row exists for the
// request -- a request canceled before any offer was ever accepted never
// had a payment created in the first place, and that's not an error here.
func (s *Service) onRequestCanceled(ctx context.Context, event broker.Event) error {
	requestID, ok := payloadInt64(event.Payload, "request_id")
	if !ok {
		return fmt.Errorf("request.canceled event missing request_id")
	}
	if _, err := s.repo.GetByRequestID(ctx, requestID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("look up payment for request %d: %w", requestID, err)
	}
	if _, err := s.repo.Refund(ctx, requestID); err != nil {
		return fmt.Errorf("refund payment for request %d: %w", requestID, err)
	}
	return nil
}

// GetForRequest is restricted to the two participants of the request --
// checked via RequestPort, since payments carries no client_id/master_id of
// its own.
func (s *Service) GetForRequest(ctx context.Context, requestID, userID int64) (PaymentDTO, error) {
	clientID, masterID, found, err := s.requests.GetParticipants(ctx, requestID)
	if err != nil {
		return PaymentDTO{}, err
	}
	if !found {
		return PaymentDTO{}, ErrRequestNotFound
	}
	if clientID != userID && (masterID == nil || *masterID != userID) {
		return PaymentDTO{}, ErrForbidden
	}

	p, err := s.repo.GetByRequestID(ctx, requestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PaymentDTO{}, ErrPaymentNotFound
		}
		return PaymentDTO{}, err
	}
	return ToPaymentDTO(p), nil
}

func (s *Service) ListAdmin(ctx context.Context, page, pageSize int) (httpx.Page[PaymentDTO], error) {
	limit, offset := int32(pageSize), httpx.Offset(page, pageSize)
	items, err := s.repo.ListAdmin(ctx, limit, offset)
	if err != nil {
		return httpx.Page[PaymentDTO]{}, err
	}
	total, err := s.repo.CountAdmin(ctx)
	if err != nil {
		return httpx.Page[PaymentDTO]{}, err
	}
	out := make([]PaymentDTO, 0, len(items))
	for _, p := range items {
		out = append(out, ToPaymentDTO(p))
	}
	return httpx.NewPage(out, page, pageSize, total), nil
}

func payloadInt64(payload map[string]any, key string) (int64, bool) {
	v, ok := payload[key].(int64)
	return v, ok
}

func payloadFloat64(payload map[string]any, key string) (float64, bool) {
	v, ok := payload[key].(float64)
	return v, ok
}
