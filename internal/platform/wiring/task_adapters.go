package wiring

import (
	"context"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/address"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/auth"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/task"
)

// TaskAddressAdapter exposes the address domain as a task.AddressStore.
type TaskAddressAdapter struct {
	addresses *address.Service
}

func NewTaskAddressAdapter(addresses *address.Service) *TaskAddressAdapter {
	return &TaskAddressAdapter{addresses: addresses}
}

func (a *TaskAddressAdapter) Exists(ctx context.Context, id, tenantID int64) (bool, error) {
	return a.addresses.Exists(ctx, id, tenantID)
}

func (a *TaskAddressAdapter) GetLabel(ctx context.Context, id, tenantID int64) (string, error) {
	return a.addresses.Label(ctx, id, tenantID)
}

// TaskUserAdapter exposes the auth domain as a task.UserStore.
type TaskUserAdapter struct {
	auth *auth.Service
}

func NewTaskUserAdapter(authService *auth.Service) *TaskUserAdapter {
	return &TaskUserAdapter{auth: authService}
}

func (a *TaskUserAdapter) GetElectrician(ctx context.Context, id, tenantID int64) (task.AssigneeSummary, error) {
	user, err := a.auth.GetTenantElectrician(ctx, id, tenantID)
	if err != nil {
		return task.AssigneeSummary{}, err
	}
	return task.AssigneeSummary{
		ID:       user.ID,
		FullName: user.FullName,
		Username: user.Username,
	}, nil
}
