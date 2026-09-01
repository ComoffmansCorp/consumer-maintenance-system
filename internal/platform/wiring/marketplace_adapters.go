package wiring

import (
	"context"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/auth"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/marketplace"
)

// MarketplaceUserAdapter exposes the auth domain as a marketplace.UserStore,
// used only to resolve display names (client/master full name) on request DTOs.
type MarketplaceUserAdapter struct {
	auth *auth.Service
}

func NewMarketplaceUserAdapter(authService *auth.Service) *MarketplaceUserAdapter {
	return &MarketplaceUserAdapter{auth: authService}
}

func (a *MarketplaceUserAdapter) GetUser(ctx context.Context, id int64) (marketplace.UserSummary, error) {
	user, err := a.auth.GetUser(ctx, id)
	if err != nil {
		return marketplace.UserSummary{}, err
	}
	return marketplace.UserSummary{ID: user.ID, FullName: user.FullName}, nil
}
