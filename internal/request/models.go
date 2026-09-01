package request

import "time"

type Status string

const (
	StatusOpen      Status = "OPEN"
	StatusAssigned  Status = "ASSIGNED"
	StatusCompleted Status = "COMPLETED"
	StatusCanceled  Status = "CANCELED"
)

type OfferStatus string

const (
	OfferPending   OfferStatus = "PENDING"
	OfferAccepted  OfferStatus = "ACCEPTED"
	OfferRejected  OfferStatus = "REJECTED"
	OfferWithdrawn OfferStatus = "WITHDRAWN"
)

// ServiceRequest is a client's заявка. Platform-level, not tenant-scoped:
// any client posts a request against the shared catalog, any specialized
// master can bid on it from the open pool.
type ServiceRequest struct {
	ID           int64
	ClientID     int64
	ServiceID    int64
	Description  string
	AddressText  string
	Latitude     *float64
	Longitude    *float64
	Status       Status
	MasterID     *int64
	AgreedPrice  *float64
	CancelReason string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Offer struct {
	ID        int64
	RequestID int64
	MasterID  int64
	Price     float64
	Comment   string
	Status    OfferStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type StatusHistoryEntry struct {
	ID         int64
	RequestID  int64
	FromStatus string
	ToStatus   string
	ChangedBy  int64
	Comment    string
	CreatedAt  time.Time
}

type Favorite struct {
	ClientID  int64
	MasterID  int64
	CreatedAt time.Time
}
