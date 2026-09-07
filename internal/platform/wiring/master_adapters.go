package wiring

import (
	"context"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/master"
)

// MasterAdapter exposes the master domain to its consumers (request,
// review) as their respective port. One concrete adapter satisfies both
// interfaces structurally.
type MasterAdapter struct {
	master *master.Service
}

func NewMasterAdapter(masterService *master.Service) *MasterAdapter {
	return &MasterAdapter{master: masterService}
}

// HasSpecialization satisfies request.SpecializationPort.
func (a *MasterAdapter) HasSpecialization(ctx context.Context, masterUserID, serviceID int64) (bool, error) {
	return a.master.HasSpecialization(ctx, masterUserID, serviceID)
}

// GetAvatarURL also satisfies request.SpecializationPort -- request enriches
// each OfferDTO with the bidding master's avatar the same way it already
// enriches RequestDTO.ServiceName via CatalogPort.
func (a *MasterAdapter) GetAvatarURL(ctx context.Context, masterUserID int64) (*string, error) {
	return a.master.GetAvatarURL(ctx, masterUserID)
}

// RecordReview satisfies review.MasterPort.
func (a *MasterAdapter) RecordReview(ctx context.Context, masterUserID int64, rating int) error {
	return a.master.RecordReview(ctx, masterUserID, rating)
}
