package wiring

import (
	"context"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/act"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/address"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/auth"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/consumer"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/meter"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/organization"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/photo"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/task"
)

// ActTaskAdapter exposes the task domain as an act.TaskStore.
type ActTaskAdapter struct {
	tasks *task.Service
}

func NewActTaskAdapter(tasks *task.Service) *ActTaskAdapter {
	return &ActTaskAdapter{tasks: tasks}
}

func (a *ActTaskAdapter) GetTask(ctx context.Context, id, tenantID int64) (act.TaskInfo, error) {
	t, err := a.tasks.GetTask(ctx, id, tenantID)
	if err != nil {
		return act.TaskInfo{}, err
	}
	return act.TaskInfo{
		ID:         t.ID,
		Type:       string(t.Type),
		Status:     string(t.Status),
		TenantID:   t.TenantID,
		AddressID:  t.AddressID,
		AssigneeID: t.AssigneeID,
	}, nil
}

// ActAddressAdapter exposes the address domain as an act.AddressStore.
type ActAddressAdapter struct {
	addresses *address.Service
}

func NewActAddressAdapter(addresses *address.Service) *ActAddressAdapter {
	return &ActAddressAdapter{addresses: addresses}
}

func (a *ActAddressAdapter) Get(ctx context.Context, id, tenantID int64) (act.AddressInfo, error) {
	label, err := a.addresses.Label(ctx, id, tenantID)
	if err != nil {
		return act.AddressInfo{}, err
	}
	return act.AddressInfo{ID: id, Label: label}, nil
}

// ActConsumerAdapter exposes the consumer domain as an act.ConsumerStore.
type ActConsumerAdapter struct {
	consumers *consumer.Service
}

func NewActConsumerAdapter(consumers *consumer.Service) *ActConsumerAdapter {
	return &ActConsumerAdapter{consumers: consumers}
}

func (a *ActConsumerAdapter) Get(ctx context.Context, id, tenantID int64) (act.ConsumerInfo, error) {
	c, err := a.consumers.Summary(ctx, id, tenantID)
	if err != nil {
		return act.ConsumerInfo{}, err
	}
	return act.ConsumerInfo{ID: c.ID, Name: c.Name, Type: string(c.Type)}, nil
}

// ActMeterAdapter exposes the meter domain as an act.MeterLister.
type ActMeterAdapter struct {
	meters *meter.Service
}

func NewActMeterAdapter(meters *meter.Service) *ActMeterAdapter {
	return &ActMeterAdapter{meters: meters}
}

func (a *ActMeterAdapter) ListByAct(ctx context.Context, actID, tenantID int64) ([]act.MeterInfo, error) {
	meters, err := a.meters.ListByAct(ctx, actID, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]act.MeterInfo, 0, len(meters))
	for _, m := range meters {
		out = append(out, act.MeterInfo{
			ID:                  m.ID,
			Type:                string(m.Type),
			SerialNumber:        m.SerialNumber,
			ManufactureYear:     m.ManufactureYear,
			VerificationDate:    httpx.DateToTime(m.VerificationDate),
			SealState:           string(m.SealState),
			TransformationRatio: m.TransformationRatio,
		})
	}
	return out, nil
}

// ActPhotoAdapter exposes the photo domain as an act.PhotoLister.
type ActPhotoAdapter struct {
	photos *photo.Service
}

func NewActPhotoAdapter(photos *photo.Service) *ActPhotoAdapter {
	return &ActPhotoAdapter{photos: photos}
}

func (a *ActPhotoAdapter) ListForInspectionAct(ctx context.Context, actID int64) ([]act.PhotoInfo, error) {
	summaries, err := a.photos.ListForInspectionAct(ctx, actID)
	if err != nil {
		return nil, err
	}
	return toActPhotoInfos(summaries), nil
}

func (a *ActPhotoAdapter) ListForReplacementAct(ctx context.Context, actID int64) ([]act.PhotoInfo, error) {
	summaries, err := a.photos.ListForReplacementAct(ctx, actID)
	if err != nil {
		return nil, err
	}
	return toActPhotoInfos(summaries), nil
}

func toActPhotoInfos(summaries []photo.Summary) []act.PhotoInfo {
	out := make([]act.PhotoInfo, 0, len(summaries))
	for _, p := range summaries {
		out = append(out, act.PhotoInfo{
			ID:               p.ID,
			OriginalFilename: p.OriginalFilename,
			Note:             p.Note,
			FilePath:         p.FilePath,
			CreatedAt:        p.CreatedAt,
		})
	}
	return out
}

// ActUserAdapter exposes the auth domain as an act.UserStore.
type ActUserAdapter struct {
	auth *auth.Service
}

func NewActUserAdapter(authService *auth.Service) *ActUserAdapter {
	return &ActUserAdapter{auth: authService}
}

func (a *ActUserAdapter) Get(ctx context.Context, id int64) (act.InspectorInfo, error) {
	u, err := a.auth.GetUser(ctx, id)
	if err != nil {
		return act.InspectorInfo{}, err
	}
	return act.InspectorInfo{ID: u.ID, FullName: u.FullName}, nil
}

// ActTenantAdapter exposes the organization domain as an act.TenantStore,
// for the PDF letterhead.
type ActTenantAdapter struct {
	org *organization.Repository
}

func NewActTenantAdapter(org *organization.Repository) *ActTenantAdapter {
	return &ActTenantAdapter{org: org}
}

func (a *ActTenantAdapter) Get(ctx context.Context, id int64) (act.TenantInfo, error) {
	t, err := a.org.GetByID(ctx, id)
	if err != nil {
		return act.TenantInfo{}, err
	}
	return act.TenantInfo{Name: t.Name, Code: t.Code}, nil
}
