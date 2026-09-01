package chat

import "context"

// RequestPort confirms who may read/write a request's chat thread: the
// client and the assigned master, and only once the request has actually
// been assigned (master_id IS NOT NULL) -- there's no one to chat with on a
// still-OPEN request. found=false covers "no such request" without leaking
// request's error sentinels across the domain boundary. Satisfied by
// request.Service.
type RequestPort interface {
	GetParticipants(ctx context.Context, requestID int64) (clientID int64, masterID *int64, found bool, err error)
}
