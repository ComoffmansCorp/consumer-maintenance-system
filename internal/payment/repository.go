package payment

import (
	"context"
	"fmt"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/payment/db"
	platformdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/db"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{queries: queries}
}

func (r *Repository) queriesFor(ctx context.Context) *db.Queries {
	if tx, ok := platformdb.TxFromContext(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *Repository) Create(ctx context.Context, requestID int64, amount, platformFee float64) (Payment, error) {
	row, err := r.queriesFor(ctx).CreatePayment(ctx, db.CreatePaymentParams{
		RequestID: requestID, Amount: amount, PlatformFee: platformFee,
	})
	if err != nil {
		return Payment{}, fmt.Errorf("create payment: %w", err)
	}
	return toPayment(row), nil
}

func (r *Repository) GetByRequestID(ctx context.Context, requestID int64) (Payment, error) {
	row, err := r.queriesFor(ctx).GetPaymentByRequestID(ctx, requestID)
	if err != nil {
		return Payment{}, fmt.Errorf("get payment by request id: %w", err)
	}
	return toPayment(row), nil
}

func (r *Repository) Release(ctx context.Context, requestID int64) (Payment, error) {
	row, err := r.queriesFor(ctx).ReleasePayment(ctx, requestID)
	if err != nil {
		return Payment{}, fmt.Errorf("release payment: %w", err)
	}
	return toPayment(row), nil
}

func (r *Repository) Refund(ctx context.Context, requestID int64) (Payment, error) {
	row, err := r.queriesFor(ctx).RefundPayment(ctx, requestID)
	if err != nil {
		return Payment{}, fmt.Errorf("refund payment: %w", err)
	}
	return toPayment(row), nil
}

func (r *Repository) ListAdmin(ctx context.Context, limit, offset int32) ([]Payment, error) {
	rows, err := r.queriesFor(ctx).ListPaymentsAdmin(ctx, db.ListPaymentsAdminParams{PageLimit: limit, PageOffset: offset})
	if err != nil {
		return nil, fmt.Errorf("list payments admin: %w", err)
	}
	out := make([]Payment, 0, len(rows))
	for _, row := range rows {
		out = append(out, toPayment(row))
	}
	return out, nil
}

func (r *Repository) CountAdmin(ctx context.Context) (int64, error) {
	count, err := r.queriesFor(ctx).CountPaymentsAdmin(ctx)
	if err != nil {
		return 0, fmt.Errorf("count payments admin: %w", err)
	}
	return count, nil
}

func toPayment(row db.Payment) Payment {
	return Payment{
		ID:          row.ID,
		RequestID:   row.RequestID,
		Amount:      row.Amount,
		PlatformFee: row.PlatformFee,
		Status:      Status(row.Status),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
