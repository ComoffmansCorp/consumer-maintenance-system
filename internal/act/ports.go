package act

import (
	"context"
	"time"
)

// TaskInfo is the projection the act domain needs from task: which tenant
// and address the act belongs to, whether the caller is the assignee, and
// whether the task is in a state where the act may still be edited.
type TaskInfo struct {
	ID         int64
	Type       string
	Status     string
	TenantID   int64
	AddressID  int64
	AssigneeID *int64
}

type TaskStore interface {
	GetTask(ctx context.Context, id, tenantID int64) (TaskInfo, error)
}

type AddressInfo struct {
	ID    int64
	Label string
}

type AddressStore interface {
	Get(ctx context.Context, id, tenantID int64) (AddressInfo, error)
}

type ConsumerInfo struct {
	ID   int64
	Name string
	Type string
}

type ConsumerStore interface {
	Get(ctx context.Context, id, tenantID int64) (ConsumerInfo, error)
}

type MeterInfo struct {
	ID                  int64
	Type                string
	SerialNumber        string
	ManufactureYear     *int32
	VerificationDate    *time.Time
	SealState           string
	TransformationRatio *int32
}

// MeterLister is satisfied by the meter domain so the act PDF/detail view
// can embed the meters recorded during an inspection.
type MeterLister interface {
	ListByAct(ctx context.Context, actID, tenantID int64) ([]MeterInfo, error)
}

type PhotoInfo struct {
	ID               int64
	OriginalFilename string
	Note             string
	// FilePath is the absolute path on disk, used to embed the photo in the
	// generated PDF. Empty if the file is missing or not an embeddable image.
	FilePath  string
	CreatedAt time.Time
}

// PhotoLister is satisfied by the photo domain.
type PhotoLister interface {
	ListForInspectionAct(ctx context.Context, actID int64) ([]PhotoInfo, error)
	ListForReplacementAct(ctx context.Context, actID int64) ([]PhotoInfo, error)
}

type InspectorInfo struct {
	ID       int64
	FullName string
}

// UserStore is satisfied by the auth domain.
type UserStore interface {
	Get(ctx context.Context, id int64) (InspectorInfo, error)
}

type TenantInfo struct {
	Name string
	Code string
}

// TenantStore is satisfied by the organization domain, for the PDF letterhead.
type TenantStore interface {
	Get(ctx context.Context, id int64) (TenantInfo, error)
}
