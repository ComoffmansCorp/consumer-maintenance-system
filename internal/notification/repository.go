package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/notification/db"
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

func (r *Repository) Create(ctx context.Context, tenantID, userID int64, notifType, title, message string, payload map[string]any) (Notification, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return Notification{}, fmt.Errorf("marshal payload: %w", err)
	}
	row, err := r.queriesFor(ctx).CreateNotification(ctx, db.CreateNotificationParams{
		TenantID: tenantID,
		UserID:   userID,
		Type:     notifType,
		Title:    title,
		Message:  message,
		Payload:  payloadJSON,
	})
	if err != nil {
		return Notification{}, fmt.Errorf("create notification: %w", err)
	}
	return toNotification(row), nil
}

func (r *Repository) List(ctx context.Context, tenantID, userID int64, unreadOnly bool, limit, offset int32) ([]Notification, error) {
	rows, err := r.queriesFor(ctx).ListNotificationsByUser(ctx, db.ListNotificationsByUserParams{
		TenantID:   tenantID,
		UserID:     userID,
		UnreadOnly: unreadOnly,
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	out := make([]Notification, 0, len(rows))
	for _, row := range rows {
		out = append(out, toNotification(row))
	}
	return out, nil
}

func (r *Repository) Count(ctx context.Context, tenantID, userID int64, unreadOnly bool) (int64, error) {
	count, err := r.queriesFor(ctx).CountNotificationsByUser(ctx, db.CountNotificationsByUserParams{
		TenantID:   tenantID,
		UserID:     userID,
		UnreadOnly: unreadOnly,
	})
	if err != nil {
		return 0, fmt.Errorf("count notifications: %w", err)
	}
	return count, nil
}

func (r *Repository) UnreadCount(ctx context.Context, tenantID, userID int64) (int64, error) {
	count, err := r.queriesFor(ctx).CountUnreadByUser(ctx, db.CountUnreadByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return 0, fmt.Errorf("count unread notifications: %w", err)
	}
	return count, nil
}

func (r *Repository) MarkRead(ctx context.Context, id, tenantID, userID int64) (Notification, error) {
	row, err := r.queriesFor(ctx).MarkNotificationRead(ctx, db.MarkNotificationReadParams{
		ID:       id,
		TenantID: tenantID,
		UserID:   userID,
	})
	if err != nil {
		return Notification{}, fmt.Errorf("mark notification read: %w", err)
	}
	return toNotification(row), nil
}

func (r *Repository) MarkAllRead(ctx context.Context, tenantID, userID int64) error {
	if err := r.queriesFor(ctx).MarkAllNotificationsRead(ctx, db.MarkAllNotificationsReadParams{
		TenantID: tenantID,
		UserID:   userID,
	}); err != nil {
		return fmt.Errorf("mark all notifications read: %w", err)
	}
	return nil
}

func toNotification(row db.Notification) Notification {
	var payload map[string]any
	_ = json.Unmarshal(row.Payload, &payload)
	var readAt *pgtype.Timestamptz
	if row.ReadAt.Valid {
		readAt = &row.ReadAt
	}
	n := Notification{
		ID:        row.ID,
		TenantID:  row.TenantID,
		UserID:    row.UserID,
		Type:      row.Type,
		Title:     row.Title,
		Message:   row.Message,
		Payload:   payload,
		CreatedAt: row.CreatedAt.Time,
	}
	if readAt != nil {
		t := readAt.Time
		n.ReadAt = &t
	}
	return n
}
