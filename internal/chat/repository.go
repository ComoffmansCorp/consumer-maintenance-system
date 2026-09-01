package chat

import (
	"context"
	"fmt"
	"time"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/chat/db"
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

func (r *Repository) Create(ctx context.Context, requestID, senderID int64, text string) (Message, error) {
	row, err := r.queriesFor(ctx).CreateMessage(ctx, db.CreateMessageParams{
		RequestID: requestID, SenderID: senderID, Text: text,
	})
	if err != nil {
		return Message{}, fmt.Errorf("create message: %w", err)
	}
	return toMessage(row), nil
}

func (r *Repository) ListSince(ctx context.Context, requestID, sinceID int64, limit int32) ([]Message, error) {
	rows, err := r.queriesFor(ctx).ListMessagesByRequestSince(ctx, db.ListMessagesByRequestSinceParams{
		RequestID: requestID, SinceID: sinceID, PageLimit: limit,
	})
	if err != nil {
		return nil, fmt.Errorf("list messages since: %w", err)
	}
	out := make([]Message, 0, len(rows))
	for _, row := range rows {
		out = append(out, toMessage(row))
	}
	return out, nil
}

func toMessage(row db.Message) Message {
	var readAt *time.Time
	if row.ReadAt.Valid {
		readAt = &row.ReadAt.Time
	}
	return Message{
		ID:        row.ID,
		RequestID: row.RequestID,
		SenderID:  row.SenderID,
		Text:      row.Text,
		CreatedAt: row.CreatedAt.Time,
		ReadAt:    readAt,
	}
}
