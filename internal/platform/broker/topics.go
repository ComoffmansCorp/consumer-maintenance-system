package broker

// Well-known event types shared between publishing and subscribing domains.
// Keeping them here (rather than in a domain package) lets notification
// subscribe without importing task/act as concrete dependencies.
const (
	// Payload: task_id, address_id, assignee_id, assignee_name, task_type.
	EventTaskAssigned = "task.assigned"
	// Payload: task_id, status, assignee_id, task_type.
	EventTaskStatusChanged = "task.status_changed"
)
