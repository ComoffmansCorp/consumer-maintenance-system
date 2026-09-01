package broker

// Well-known event types shared between publishing and subscribing domains.
// Keeping them here (rather than in a domain package) lets payment subscribe
// without importing request as a concrete dependency.
const (
	// Payload: request_id, master_id, price.
	EventRequestAssigned = "request.assigned"
	// Payload: request_id.
	EventRequestCompleted = "request.completed"
	// Payload: request_id.
	EventRequestCanceled = "request.canceled"
)
