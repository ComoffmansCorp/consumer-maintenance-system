package task

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/task/db"
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

func (r *Repository) Create(ctx context.Context, tenantID int64, typ Type, addressID int64, dueDate *time.Time, assigneeID *int64) (Task, error) {
	row, err := r.queriesFor(ctx).CreateTask(ctx, db.CreateTaskParams{
		Type:       string(typ),
		TenantID:   tenantID,
		AddressID:  addressID,
		DueDate:    dateOrNull(dueDate),
		AssigneeID: int8OrNull(assigneeID),
	})
	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}
	return toTask(row), nil
}

func (r *Repository) GetByID(ctx context.Context, id, tenantID int64) (Task, error) {
	row, err := r.queriesFor(ctx).GetTaskByID(ctx, db.GetTaskByIDParams{ID: id, TenantID: tenantID})
	if err != nil {
		return Task{}, fmt.Errorf("get task by id: %w", err)
	}
	return toTask(row), nil
}

func (r *Repository) List(ctx context.Context, tenantID int64, status Status, typ Type, assigneeID int64, limit, offset int32) ([]Task, error) {
	rows, err := r.queriesFor(ctx).ListTasks(ctx, db.ListTasksParams{
		TenantID:   tenantID,
		Status:     string(status),
		TaskType:   string(typ),
		AssigneeID: assigneeID,
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	out := make([]Task, 0, len(rows))
	for _, row := range rows {
		out = append(out, toTask(row))
	}
	return out, nil
}

func (r *Repository) Count(ctx context.Context, tenantID int64, status Status, typ Type, assigneeID int64) (int64, error) {
	count, err := r.queriesFor(ctx).CountTasks(ctx, db.CountTasksParams{
		TenantID:   tenantID,
		Status:     string(status),
		TaskType:   string(typ),
		AssigneeID: assigneeID,
	})
	if err != nil {
		return 0, fmt.Errorf("count tasks: %w", err)
	}
	return count, nil
}

func (r *Repository) Assign(ctx context.Context, id, tenantID, assigneeID int64) (Task, error) {
	row, err := r.queriesFor(ctx).AssignTask(ctx, db.AssignTaskParams{
		AssigneeID: pgtype.Int8{Int64: assigneeID, Valid: true},
		ID:         id,
		TenantID:   tenantID,
	})
	if err != nil {
		return Task{}, fmt.Errorf("assign task: %w", err)
	}
	return toTask(row), nil
}

func (r *Repository) Start(ctx context.Context, id, tenantID int64) (Task, error) {
	row, err := r.queriesFor(ctx).StartTask(ctx, db.StartTaskParams{ID: id, TenantID: tenantID})
	if err != nil {
		return Task{}, fmt.Errorf("start task: %w", err)
	}
	return toTask(row), nil
}

func (r *Repository) Complete(ctx context.Context, id, tenantID int64) (Task, error) {
	row, err := r.queriesFor(ctx).CompleteTask(ctx, db.CompleteTaskParams{ID: id, TenantID: tenantID})
	if err != nil {
		return Task{}, fmt.Errorf("complete task: %w", err)
	}
	return toTask(row), nil
}

func (r *Repository) Cancel(ctx context.Context, id, tenantID int64, reason string) (Task, error) {
	row, err := r.queriesFor(ctx).CancelTask(ctx, db.CancelTaskParams{
		CancelReason: textOrNull(reason),
		ID:           id,
		TenantID:     tenantID,
	})
	if err != nil {
		return Task{}, fmt.Errorf("cancel task: %w", err)
	}
	return toTask(row), nil
}

func dateOrNull(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func int8OrNull(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}

func textOrNull(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func toTask(row db.Task) Task {
	var dueDate *time.Time
	if row.DueDate.Valid {
		dueDate = &row.DueDate.Time
	}
	var assigneeID *int64
	if row.AssigneeID.Valid {
		assigneeID = &row.AssigneeID.Int64
	}
	var completedAt *time.Time
	if row.CompletedAt.Valid {
		completedAt = &row.CompletedAt.Time
	}
	var canceledAt *time.Time
	if row.CanceledAt.Valid {
		canceledAt = &row.CanceledAt.Time
	}
	return Task{
		ID:           row.ID,
		Type:         Type(row.Type),
		TenantID:     row.TenantID,
		AddressID:    row.AddressID,
		Status:       Status(row.Status),
		DueDate:      dueDate,
		AssigneeID:   assigneeID,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		CompletedAt:  completedAt,
		CanceledAt:   canceledAt,
		CancelReason: row.CancelReason.String,
	}
}
