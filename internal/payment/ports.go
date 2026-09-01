package payment

import "context"

// RequestPort resolves who a request belongs to, to authorize
// GET /requests/{id}/payment to the two participants only. found=false
// covers "no such request" without leaking request's error sentinels
// across the domain boundary. Satisfied by request.Service.
type RequestPort interface {
	GetParticipants(ctx context.Context, requestID int64) (clientID int64, masterID *int64, found bool, err error)
}
