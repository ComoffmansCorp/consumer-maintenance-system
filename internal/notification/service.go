package notification

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/broker"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

var ErrNotificationNotFound = errors.New("notification not found")

type Service struct {
	repo   *Repository
	admins AdminLister
}

func NewService(repo *Repository, admins AdminLister) *Service {
	return &Service{repo: repo, admins: admins}
}

// RegisterHandlers subscribes to the cross-domain events that should turn
// into in-app notifications. Called once from the composition root.
func (s *Service) RegisterHandlers(bus *broker.Bus) {
	bus.Subscribe(broker.EventTaskAssigned, s.onTaskAssigned)
	bus.Subscribe(broker.EventTaskStatusChanged, s.onTaskStatusChanged)
}

func (s *Service) onTaskAssigned(ctx context.Context, event broker.Event) error {
	assigneeID, _ := event.Payload["assignee_id"].(int64)
	taskID, _ := event.Payload["task_id"].(int64)
	if assigneeID == 0 {
		return nil
	}
	title := "Новый наряд"
	message := fmt.Sprintf("Вам назначен наряд №%d", taskID)
	_, err := s.repo.Create(ctx, event.TenantID, assigneeID, TypeTaskAssigned, title, message, event.Payload)
	return err
}

func (s *Service) onTaskStatusChanged(ctx context.Context, event broker.Event) error {
	status, _ := event.Payload["status"].(string)
	if status != "COMPLETED" && status != "CANCELED" {
		return nil
	}
	taskID, _ := event.Payload["task_id"].(int64)

	adminIDs, err := s.admins.ListTenantAdminIDs(ctx, event.TenantID)
	if err != nil {
		return fmt.Errorf("list tenant admins: %w", err)
	}

	title := "Наряд завершён"
	message := fmt.Sprintf("Наряд №%d выполнен", taskID)
	if status == "CANCELED" {
		title = "Наряд отменён"
		message = fmt.Sprintf("Наряд №%d отменён", taskID)
	}

	for _, adminID := range adminIDs {
		if _, err := s.repo.Create(ctx, event.TenantID, adminID, TypeTaskStatusChanged, title, message, event.Payload); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) List(ctx context.Context, tenantID, userID int64, unreadOnly bool, page, pageSize int) (httpx.Page[DTO], error) {
	limit := int32(pageSize)
	offset := httpx.Offset(page, pageSize)

	items, err := s.repo.List(ctx, tenantID, userID, unreadOnly, limit, offset)
	if err != nil {
		return httpx.Page[DTO]{}, err
	}
	total, err := s.repo.Count(ctx, tenantID, userID, unreadOnly)
	if err != nil {
		return httpx.Page[DTO]{}, err
	}

	dtos := make([]DTO, 0, len(items))
	for _, n := range items {
		dtos = append(dtos, ToDTO(n))
	}
	return httpx.NewPage(dtos, page, pageSize, total), nil
}

func (s *Service) UnreadCount(ctx context.Context, tenantID, userID int64) (int64, error) {
	return s.repo.UnreadCount(ctx, tenantID, userID)
}

func (s *Service) MarkRead(ctx context.Context, id, tenantID, userID int64) (DTO, error) {
	n, err := s.repo.MarkRead(ctx, id, tenantID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO{}, ErrNotificationNotFound
		}
		return DTO{}, err
	}
	return ToDTO(n), nil
}

func (s *Service) MarkAllRead(ctx context.Context, tenantID, userID int64) error {
	return s.repo.MarkAllRead(ctx, tenantID, userID)
}
