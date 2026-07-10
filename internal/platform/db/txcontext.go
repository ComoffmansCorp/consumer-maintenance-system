package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type txKey struct{}

func ContextWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

type TxManager struct {
	pool *Pool
}

func NewTxManager(pool *Pool) *TxManager {
	return &TxManager{pool: pool}
}

func (m *TxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.pool.WithTx(ctx, func(tx pgx.Tx) error {
		return fn(ContextWithTx(ctx, tx))
	})
}
