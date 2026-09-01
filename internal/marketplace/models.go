package marketplace

import "time"

type RequestStatus string

const (
	RequestStatusOpen       RequestStatus = "OPEN"
	RequestStatusInProgress RequestStatus = "IN_PROGRESS"
	RequestStatusCompleted  RequestStatus = "COMPLETED"
	RequestStatusCanceled   RequestStatus = "CANCELED"
)

type Category struct {
	ID        int64
	Name      string
	Active    bool
	CreatedAt time.Time
}

// Offering is a concrete service a client can order (e.g. "Замена счётчика
// воды") within a Category. Named Offering rather than Service to avoid
// colliding with this package's own business-logic Service struct below --
// "service" the architectural layer and "service" the marketplace noun are
// unrelated concepts that happen to share an English word.
type Offering struct {
	ID          int64
	CategoryID  int64
	Name        string
	Description string
	Active      bool
	CreatedAt   time.Time
}

type MasterProfile struct {
	UserID    int64
	City      string
	Bio       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Request is a client's заявка -- platform-level, not tenant-scoped. Distinct
// from the task domain's наряд: same shape of idea (work to be done), but a
// different flow (self-service claim from an open pool vs. dispatcher
// assignment) and a broader, admin-curated service catalog instead of the
// fixed INSPECTION/REPLACEMENT pair.
type Request struct {
	ID           int64
	ClientID     int64
	OfferingID   int64
	Description  string
	AddressText  string
	Latitude     *float64
	Longitude    *float64
	Status       RequestStatus
	MasterID     *int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ClaimedAt    *time.Time
	CompletedAt  *time.Time
	CanceledAt   *time.Time
	CancelReason string
}
