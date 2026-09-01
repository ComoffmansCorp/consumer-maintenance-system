package review

import "context"

// RequestPort validates that a review can be left: the request must be
// COMPLETED and belong to the reviewing client. Returns the assigned master
// to attribute the review to; eligible=false covers "not found", "not
// yours" and "not completed yet" alike, without leaking request's internal
// error sentinels across the domain boundary. Satisfied by request.Service.
type RequestPort interface {
	GetCompletedForReview(ctx context.Context, requestID, clientID int64) (masterID int64, eligible bool, err error)
}

// MasterPort folds a new rating into the master's running average.
// Satisfied by master.Service.
type MasterPort interface {
	RecordReview(ctx context.Context, masterUserID int64, rating int) error
}
