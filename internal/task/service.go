package task

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/broker"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

var (
	ErrTaskNotFound         = errors.New("task not found")
	ErrInvalidType          = errors.New("invalid task type")
	ErrAddressNotFound      = errors.New("address not found")
	ErrAssigneeNotFound     = errors.New("assignee not found or not an electrician")
	ErrInvalidTransition    = errors.New("invalid task status transition")
	ErrNotAssignee          = errors.New("only the assigned inspector can perform this action")
	ErrActNotFilled         = errors.New("act must be filled in before completing the task")
	ErrCancelReasonRequired = errors.New("cancel reason is required")
)

type Service struct {
	repo      *Repository
	addresses AddressStore
	users     UserStore
	acts      ActStore
	events    *broker.Bus
}

func NewService(repo *Repository, addresses AddressStore, users UserStore, acts ActStore, events *broker.Bus) *Service {
	return &Service{repo: repo, addresses: addresses, users: users, acts: acts, events: events}
}

func (s *Service) Create(ctx context.Context, tenantID int64, req CreateRequest) (DTO, error) {
	if !req.Type.Valid() {
		return DTO{}, ErrInvalidType
	}
	exists, err := s.addresses.Exists(ctx, req.AddressID, tenantID)
	if err != nil {
		return DTO{}, fmt.Errorf("check address: %w", err)
	}
	if !exists {
		return DTO{}, ErrAddressNotFound
	}

	var assignee *AssigneeSummary
	if req.AssigneeID != nil {
		summary, err := s.users.GetElectrician(ctx, *req.AssigneeID, tenantID)
		if err != nil {
			return DTO{}, ErrAssigneeNotFound
		}
		assignee = &summary
	}

	t, err := s.repo.Create(ctx, tenantID, req.Type, req.AddressID, httpx.DateToTime(req.DueDate), req.AssigneeID)
	if err != nil {
		return DTO{}, fmt.Errorf("create task: %w", err)
	}
	if assignee != nil {
		s.publishAssigned(ctx, t, *assignee)
	}
	return s.enrich(ctx, t), nil
}

func (s *Service) Get(ctx context.Context, id, tenantID int64) (DTO, error) {
	t, err := s.getOrNotFound(ctx, id, tenantID)
	if err != nil {
		return DTO{}, err
	}
	return s.enrich(ctx, t), nil
}

func (s *Service) List(ctx context.Context, tenantID int64, status Status, typ Type, assigneeID int64, page, pageSize int) (httpx.Page[DTO], error) {
	limit := int32(pageSize)
	offset := httpx.Offset(page, pageSize)

	tasks, err := s.repo.List(ctx, tenantID, status, typ, assigneeID, limit, offset)
	if err != nil {
		return httpx.Page[DTO]{}, err
	}
	total, err := s.repo.Count(ctx, tenantID, status, typ, assigneeID)
	if err != nil {
		return httpx.Page[DTO]{}, err
	}

	items := make([]DTO, 0, len(tasks))
	for _, t := range tasks {
		items = append(items, s.enrich(ctx, t))
	}
	return httpx.NewPage(items, page, pageSize, total), nil
}

func (s *Service) Assign(ctx context.Context, id, tenantID, assigneeID int64) (DTO, error) {
	summary, err := s.users.GetElectrician(ctx, assigneeID, tenantID)
	if err != nil {
		return DTO{}, ErrAssigneeNotFound
	}

	current, err := s.getOrNotFound(ctx, id, tenantID)
	if err != nil {
		return DTO{}, err
	}
	if current.Status != StatusPending {
		return DTO{}, ErrInvalidTransition
	}

	t, err := s.repo.Assign(ctx, id, tenantID, assigneeID)
	if err != nil {
		return DTO{}, fmt.Errorf("assign task: %w", err)
	}
	s.publishAssigned(ctx, t, summary)
	return s.enrich(ctx, t), nil
}

func (s *Service) Start(ctx context.Context, id, tenantID, userID int64) (DTO, error) {
	current, err := s.getOrNotFound(ctx, id, tenantID)
	if err != nil {
		return DTO{}, err
	}
	if current.AssigneeID == nil || *current.AssigneeID != userID {
		return DTO{}, ErrNotAssignee
	}
	if current.Status != StatusPending {
		return DTO{}, ErrInvalidTransition
	}

	t, err := s.repo.Start(ctx, id, tenantID)
	if err != nil {
		return DTO{}, fmt.Errorf("start task: %w", err)
	}
	s.publishStatusChanged(ctx, t)
	return s.enrich(ctx, t), nil
}

func (s *Service) Complete(ctx context.Context, id, tenantID, userID int64) (DTO, error) {
	current, err := s.getOrNotFound(ctx, id, tenantID)
	if err != nil {
		return DTO{}, err
	}
	if current.AssigneeID == nil || *current.AssigneeID != userID {
		return DTO{}, ErrNotAssignee
	}
	if current.Status != StatusInProgress {
		return DTO{}, ErrInvalidTransition
	}

	hasAct, err := s.acts.HasActForTask(ctx, id, string(current.Type))
	if err != nil {
		return DTO{}, fmt.Errorf("check act: %w", err)
	}
	if !hasAct {
		return DTO{}, ErrActNotFilled
	}

	t, err := s.repo.Complete(ctx, id, tenantID)
	if err != nil {
		return DTO{}, fmt.Errorf("complete task: %w", err)
	}
	s.publishStatusChanged(ctx, t)
	return s.enrich(ctx, t), nil
}

func (s *Service) Cancel(ctx context.Context, id, tenantID int64, reason string) (DTO, error) {
	if strings.TrimSpace(reason) == "" {
		return DTO{}, ErrCancelReasonRequired
	}
	current, err := s.getOrNotFound(ctx, id, tenantID)
	if err != nil {
		return DTO{}, err
	}
	if current.Status == StatusCompleted || current.Status == StatusCanceled {
		return DTO{}, ErrInvalidTransition
	}

	t, err := s.repo.Cancel(ctx, id, tenantID, reason)
	if err != nil {
		return DTO{}, fmt.Errorf("cancel task: %w", err)
	}
	s.publishStatusChanged(ctx, t)
	return s.enrich(ctx, t), nil
}

// GetTask returns the raw domain model (including status and tenant id),
// satisfying the act domain's TaskStore port.
func (s *Service) GetTask(ctx context.Context, id, tenantID int64) (Task, error) {
	return s.getOrNotFound(ctx, id, tenantID)
}

func (s *Service) getOrNotFound(ctx context.Context, id, tenantID int64) (Task, error) {
	t, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Task{}, ErrTaskNotFound
		}
		return Task{}, err
	}
	return t, nil
}

func (s *Service) enrich(ctx context.Context, t Task) DTO {
	dto := ToDTO(t)
	if label, err := s.addresses.GetLabel(ctx, t.AddressID, t.TenantID); err == nil {
		dto.AddressLabel = label
	}
	if t.AssigneeID != nil {
		if summary, err := s.users.GetElectrician(ctx, *t.AssigneeID, t.TenantID); err == nil {
			dto.AssigneeName = summary.FullName
		}
	}
	return dto
}

func (s *Service) publishAssigned(ctx context.Context, t Task, assignee AssigneeSummary) {
	s.events.Publish(ctx, broker.Event{
		Type:     broker.EventTaskAssigned,
		TenantID: t.TenantID,
		Payload: map[string]any{
			"task_id":       t.ID,
			"task_type":     string(t.Type),
			"address_id":    t.AddressID,
			"assignee_id":   assignee.ID,
			"assignee_name": assignee.FullName,
		},
	})
}

func (s *Service) publishStatusChanged(ctx context.Context, t Task) {
	var assigneeID int64
	if t.AssigneeID != nil {
		assigneeID = *t.AssigneeID
	}
	s.events.Publish(ctx, broker.Event{
		Type:     broker.EventTaskStatusChanged,
		TenantID: t.TenantID,
		Payload: map[string]any{
			"task_id":     t.ID,
			"task_type":   string(t.Type),
			"status":      string(t.Status),
			"assignee_id": assigneeID,
		},
	})
}
