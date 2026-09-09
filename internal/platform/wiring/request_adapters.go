package wiring

import (
	"context"
	"errors"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/request"
)

// RequestAdapter exposes the request domain to its consumers (review,
// payment, chat) as their respective port. One concrete adapter satisfies
// all three interfaces structurally -- each translates request's own typed
// errors into a plain found/eligible bool so the internal error sentinels
// never leak across the domain boundary.
type RequestAdapter struct {
	requests *request.Service
}

func NewRequestAdapter(requestService *request.Service) *RequestAdapter {
	return &RequestAdapter{requests: requestService}
}

// GetCompletedForReview satisfies review.RequestPort.
func (a *RequestAdapter) GetCompletedForReview(ctx context.Context, requestID, clientID int64) (int64, bool, error) {
	masterID, err := a.requests.GetCompletedForReview(ctx, requestID, clientID)
	if err != nil {
		if isRequestIneligible(err) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return masterID, true, nil
}

// GetParticipants satisfies both payment.RequestPort and chat.RequestPort.
func (a *RequestAdapter) GetParticipants(ctx context.Context, requestID int64) (int64, *int64, bool, error) {
	clientID, masterID, err := a.requests.GetParticipants(ctx, requestID)
	if err != nil {
		if errors.Is(err, request.ErrRequestNotFound) {
			return 0, nil, false, nil
		}
		return 0, nil, false, err
	}
	return clientID, masterID, true, nil
}

// ListParticipantRequestIDs satisfies chat.RequestPort.
func (a *RequestAdapter) ListParticipantRequestIDs(ctx context.Context, userID int64) ([]int64, error) {
	return a.requests.ListParticipantRequestIDs(ctx, userID)
}

func isRequestIneligible(err error) bool {
	return errors.Is(err, request.ErrRequestNotFound) ||
		errors.Is(err, request.ErrNotOwner) ||
		errors.Is(err, request.ErrInvalidTransition)
}
