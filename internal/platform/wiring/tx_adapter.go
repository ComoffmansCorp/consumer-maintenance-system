package wiring

import (
	"context"
	"fmt"
)

// TxRunnerAdapter wraps platformdb.TxManager (or anything with the same
// shape) as a domain-local TxRunner port. Used by the request domain for
// AcceptOffer, which moves the offer, its rejected siblings, and the
// request itself in one commit.
type TxRunnerAdapter struct {
	manager interface {
		WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	}
}

func NewTxRunnerAdapter(manager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}) *TxRunnerAdapter {
	return &TxRunnerAdapter{manager: manager}
}

func (a *TxRunnerAdapter) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := a.manager.WithinTransaction(ctx, fn); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}
