package payment

import "time"

type Status string

const (
	StatusHeld     Status = "HELD"
	StatusReleased Status = "RELEASED"
	StatusRefunded Status = "REFUNDED"
)

// Payment is an escrow-style ledger row per request. Created and moved
// through its lifecycle entirely by broker events published from the
// request domain -- there is no direct write path from a handler.
type Payment struct {
	ID          int64
	RequestID   int64
	Amount      float64
	PlatformFee float64
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
