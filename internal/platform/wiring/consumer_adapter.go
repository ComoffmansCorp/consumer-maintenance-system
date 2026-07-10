package wiring

import (
	"context"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/address"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/consumer"
)

// ConsumerAdapter exposes the consumer domain as an address.ConsumerStore
// port implementation, keeping the two domains decoupled from each other's
// concrete types.
type ConsumerAdapter struct {
	consumers *consumer.Service
}

func NewConsumerAdapter(consumers *consumer.Service) *ConsumerAdapter {
	return &ConsumerAdapter{consumers: consumers}
}

func (a *ConsumerAdapter) Exists(ctx context.Context, id, tenantID int64) (bool, error) {
	return a.consumers.Exists(ctx, id, tenantID)
}

func (a *ConsumerAdapter) GetSummary(ctx context.Context, id, tenantID int64) (address.ConsumerSummary, error) {
	c, err := a.consumers.Summary(ctx, id, tenantID)
	if err != nil {
		return address.ConsumerSummary{}, err
	}
	return address.ConsumerSummary{
		ID:   c.ID,
		Name: c.Name,
		Type: string(c.Type),
	}, nil
}
