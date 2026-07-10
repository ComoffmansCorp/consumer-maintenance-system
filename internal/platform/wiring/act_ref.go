package wiring

import "context"

// ActRef breaks the construction cycle between act and the domains it
// depends on (task, meter, photo) which in turn need to call back into act
// (task.ActStore, meter.ActStore, photo.ActStore). It is created empty,
// handed to those domains' services, and populated with the real act
// service once that is constructed — before the HTTP server starts
// accepting requests, so no caller ever observes the nil state.
type ActRef struct {
	svc actService
}

// actService is the subset of act.Service this ref forwards to. Defined
// locally (rather than importing act.Service's concrete type in the field)
// so this file has no import cycle risk if act ever needs wiring helpers.
type actService interface {
	HasActForTask(ctx context.Context, taskID int64, taskType string) (bool, error)
	EnsureInspectionAct(ctx context.Context, actID, tenantID int64) error
	EnsureReplacementAct(ctx context.Context, actID, tenantID int64) error
}

func (r *ActRef) Set(svc actService) {
	r.svc = svc
}

func (r *ActRef) HasActForTask(ctx context.Context, taskID int64, taskType string) (bool, error) {
	return r.svc.HasActForTask(ctx, taskID, taskType)
}

func (r *ActRef) EnsureInspectionAct(ctx context.Context, actID, tenantID int64) error {
	return r.svc.EnsureInspectionAct(ctx, actID, tenantID)
}

func (r *ActRef) EnsureReplacementAct(ctx context.Context, actID, tenantID int64) error {
	return r.svc.EnsureReplacementAct(ctx, actID, tenantID)
}
