package wiring

import (
	"context"
	"errors"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/catalog"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/request"
)

// CatalogAdapter exposes the catalog domain to its consumers (master,
// request) as their respective CatalogPort. One concrete adapter satisfies
// both interfaces structurally -- master only needs the active check,
// request also wants the service name to enrich its DTOs.
type CatalogAdapter struct {
	catalog *catalog.Service
}

func NewCatalogAdapter(catalogService *catalog.Service) *CatalogAdapter {
	return &CatalogAdapter{catalog: catalogService}
}

// IsServiceActive satisfies master.CatalogPort.
func (a *CatalogAdapter) IsServiceActive(ctx context.Context, serviceID int64) (bool, error) {
	return a.catalog.IsServiceActive(ctx, serviceID)
}

// GetService satisfies request.CatalogPort.
func (a *CatalogAdapter) GetService(ctx context.Context, serviceID int64) (request.ServiceInfo, bool, error) {
	dto, err := a.catalog.GetServiceByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, catalog.ErrServiceNotFound) {
			return request.ServiceInfo{}, false, nil
		}
		return request.ServiceInfo{}, false, err
	}
	return request.ServiceInfo{ID: dto.ID, Name: dto.Name, Active: dto.Active}, true, nil
}
