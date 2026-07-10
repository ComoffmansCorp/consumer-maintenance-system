package act

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

var (
	ErrTaskNotFound          = errors.New("task not found")
	ErrWrongTaskType         = errors.New("task type does not match this act kind")
	ErrNotAssignee           = errors.New("only the assigned inspector can manage this act")
	ErrTaskNotInProgress     = errors.New("task must be in progress to manage its act")
	ErrActAlreadyExists      = errors.New("an act already exists for this task")
	ErrActNotFound           = errors.New("act not found")
	ErrConsumerNotFound      = errors.New("consumer not found")
	ErrInvalidInspectionType = errors.New("invalid inspection type")
	ErrAccountNumberRequired = errors.New("account number is required")
)

const (
	taskTypeInspection  = "INSPECTION"
	taskTypeReplacement = "REPLACEMENT"
	taskStatusInProgress = "IN_PROGRESS"
)

type Service struct {
	repo      *Repository
	tasks     TaskStore
	addresses AddressStore
	consumers ConsumerStore
	meters    MeterLister
	photos    PhotoLister
	users     UserStore
	tenants   TenantStore
}

func NewService(repo *Repository, tasks TaskStore, addresses AddressStore, consumers ConsumerStore, meters MeterLister, photos PhotoLister, users UserStore, tenants TenantStore) *Service {
	return &Service{
		repo:      repo,
		tasks:     tasks,
		addresses: addresses,
		consumers: consumers,
		meters:    meters,
		photos:    photos,
		users:     users,
		tenants:   tenants,
	}
}

// --- Inspection acts ---

func (s *Service) CreateInspection(ctx context.Context, tenantID, userID int64, req CreateInspectionRequest) (InspectionDTO, error) {
	t, err := s.editableTask(ctx, req.TaskID, tenantID, userID, taskTypeInspection)
	if err != nil {
		return InspectionDTO{}, err
	}

	exists, err := s.repo.ExistsInspectionByTaskID(ctx, req.TaskID)
	if err != nil {
		return InspectionDTO{}, err
	}
	if exists {
		return InspectionDTO{}, ErrActAlreadyExists
	}

	if req.InspectionType == "" {
		req.InspectionType = InspectionScheduled
	} else if !req.InspectionType.Valid() {
		return InspectionDTO{}, ErrInvalidInspectionType
	}
	if err := s.checkConsumer(ctx, tenantID, req.ConsumerID); err != nil {
		return InspectionDTO{}, err
	}

	a, err := s.repo.CreateInspection(ctx, req.TaskID, tenantID, t.AddressID, req.ConsumerID, req)
	if err != nil {
		return InspectionDTO{}, fmt.Errorf("create inspection act: %w", err)
	}
	return s.enrichInspection(ctx, a), nil
}

func (s *Service) GetInspection(ctx context.Context, id, tenantID int64) (InspectionDTO, error) {
	a, err := s.repo.GetInspectionByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InspectionDTO{}, ErrActNotFound
		}
		return InspectionDTO{}, err
	}
	return s.enrichInspection(ctx, a), nil
}

func (s *Service) GetInspectionByTask(ctx context.Context, taskID, tenantID int64) (InspectionDTO, error) {
	a, err := s.repo.GetInspectionByTaskID(ctx, taskID, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InspectionDTO{}, ErrActNotFound
		}
		return InspectionDTO{}, err
	}
	return s.enrichInspection(ctx, a), nil
}

func (s *Service) UpdateInspection(ctx context.Context, id, tenantID, userID int64, req UpdateInspectionRequest) (InspectionDTO, error) {
	current, err := s.repo.GetInspectionByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InspectionDTO{}, ErrActNotFound
		}
		return InspectionDTO{}, err
	}
	if _, err := s.editableTask(ctx, current.TaskID, tenantID, userID, taskTypeInspection); err != nil {
		return InspectionDTO{}, err
	}
	if !req.InspectionType.Valid() {
		return InspectionDTO{}, ErrInvalidInspectionType
	}
	if err := s.checkConsumer(ctx, tenantID, req.ConsumerID); err != nil {
		return InspectionDTO{}, err
	}

	a, err := s.repo.UpdateInspection(ctx, id, tenantID, req.ConsumerID, req)
	if err != nil {
		return InspectionDTO{}, fmt.Errorf("update inspection act: %w", err)
	}
	return s.enrichInspection(ctx, a), nil
}

// --- Replacement acts ---

func (s *Service) CreateReplacement(ctx context.Context, tenantID, userID int64, req CreateReplacementRequest) (ReplacementDTO, error) {
	t, err := s.editableTask(ctx, req.TaskID, tenantID, userID, taskTypeReplacement)
	if err != nil {
		return ReplacementDTO{}, err
	}
	if strings.TrimSpace(req.AccountNumber) == "" {
		return ReplacementDTO{}, ErrAccountNumberRequired
	}

	exists, err := s.repo.ExistsReplacementByTaskID(ctx, req.TaskID)
	if err != nil {
		return ReplacementDTO{}, err
	}
	if exists {
		return ReplacementDTO{}, ErrActAlreadyExists
	}

	a, err := s.repo.CreateReplacement(ctx, req.TaskID, tenantID, t.AddressID, req)
	if err != nil {
		return ReplacementDTO{}, fmt.Errorf("create replacement act: %w", err)
	}
	return s.enrichReplacement(ctx, a), nil
}

func (s *Service) GetReplacement(ctx context.Context, id, tenantID int64) (ReplacementDTO, error) {
	a, err := s.repo.GetReplacementByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ReplacementDTO{}, ErrActNotFound
		}
		return ReplacementDTO{}, err
	}
	return s.enrichReplacement(ctx, a), nil
}

func (s *Service) GetReplacementByTask(ctx context.Context, taskID, tenantID int64) (ReplacementDTO, error) {
	a, err := s.repo.GetReplacementByTaskID(ctx, taskID, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ReplacementDTO{}, ErrActNotFound
		}
		return ReplacementDTO{}, err
	}
	return s.enrichReplacement(ctx, a), nil
}

func (s *Service) UpdateReplacement(ctx context.Context, id, tenantID, userID int64, req UpdateReplacementRequest) (ReplacementDTO, error) {
	current, err := s.repo.GetReplacementByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ReplacementDTO{}, ErrActNotFound
		}
		return ReplacementDTO{}, err
	}
	if _, err := s.editableTask(ctx, current.TaskID, tenantID, userID, taskTypeReplacement); err != nil {
		return ReplacementDTO{}, err
	}
	if strings.TrimSpace(req.AccountNumber) == "" {
		return ReplacementDTO{}, ErrAccountNumberRequired
	}

	a, err := s.repo.UpdateReplacement(ctx, id, tenantID, req)
	if err != nil {
		return ReplacementDTO{}, fmt.Errorf("update replacement act: %w", err)
	}
	return s.enrichReplacement(ctx, a), nil
}

// --- PDF generation ---

func (s *Service) GenerateInspectionPDF(ctx context.Context, id, tenantID int64) ([]byte, string, error) {
	a, err := s.repo.GetInspectionByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", ErrActNotFound
		}
		return nil, "", err
	}

	doc, err := newPDFDocument()
	if err != nil {
		return nil, "", err
	}

	tenant, _ := s.tenants.Get(ctx, tenantID)
	addr, _ := s.addresses.Get(ctx, a.AddressID, tenantID)
	meters, _ := s.meters.ListByAct(ctx, a.ID, tenantID)
	photos, _ := s.photos.ListForInspectionAct(ctx, a.ID)

	var consumerName string
	if a.ConsumerID != nil {
		if c, err := s.consumers.Get(ctx, *a.ConsumerID, tenantID); err == nil {
			consumerName = c.Name
		}
	}
	var inspectorName string
	if t, err := s.tasks.GetTask(ctx, a.TaskID, tenantID); err == nil && t.AssigneeID != nil {
		if u, err := s.users.Get(ctx, *t.AssigneeID); err == nil {
			inspectorName = u.FullName
		}
	}

	doc.letterhead(tenant, "Акт осмотра прибора учёта", fmt.Sprintf("ИНС-%d", a.ID), time.Now())
	doc.row("Адрес", addr.Label)
	if consumerName != "" {
		doc.row("Потребитель", consumerName)
	}
	doc.row("Дата осмотра", formatDate(a.InspectionDate))
	doc.row("Вид осмотра", translateInspectionType(a.InspectionType))
	doc.row("Инспектор", inspectorName)
	if a.Notes != "" {
		doc.row("Примечания", a.Notes)
	}
	doc.subheading("Приборы учёта")
	doc.metersTable(meters)
	doc.photosAppendix(photos)
	doc.hline()
	doc.text(pdfFontRegular, 9, "Подпись инспектора: ____________________        Подпись потребителя: ____________________")

	return doc.bytes(), fmt.Sprintf("inspection-act-%d.pdf", a.ID), nil
}

func (s *Service) GenerateReplacementPDF(ctx context.Context, id, tenantID int64) ([]byte, string, error) {
	a, err := s.repo.GetReplacementByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", ErrActNotFound
		}
		return nil, "", err
	}

	doc, err := newPDFDocument()
	if err != nil {
		return nil, "", err
	}

	tenant, _ := s.tenants.Get(ctx, tenantID)
	addr, _ := s.addresses.Get(ctx, a.AddressID, tenantID)
	photos, _ := s.photos.ListForReplacementAct(ctx, a.ID)

	var inspectorName string
	if t, err := s.tasks.GetTask(ctx, a.TaskID, tenantID); err == nil && t.AssigneeID != nil {
		if u, err := s.users.Get(ctx, *t.AssigneeID); err == nil {
			inspectorName = u.FullName
		}
	}

	doc.letterhead(tenant, "Акт замены прибора учёта", fmt.Sprintf("ЗАМ-%d", a.ID), time.Now())
	doc.row("Адрес", addr.Label)
	doc.row("Лицевой счёт", a.AccountNumber)
	doc.row("Дата замены", formatDate(a.InstallationDate))
	doc.row("Инспектор", inspectorName)
	doc.subheading("Снятый прибор")
	doc.row("Марка", a.OldBrand)
	doc.row("Серийный номер", a.OldSerialNumber)
	doc.row("Показания", formatReadings(a.OldReadings))
	doc.subheading("Установленный прибор")
	doc.row("Марка", a.NewBrand)
	doc.row("Серийный номер", a.NewSerialNumber)
	doc.row("Показания", formatReadings(a.NewReadings))
	doc.photosAppendix(photos)
	doc.hline()
	doc.text(pdfFontRegular, 9, "Подпись инспектора: ____________________        Подпись потребителя: ____________________")

	return doc.bytes(), fmt.Sprintf("replacement-act-%d.pdf", a.ID), nil
}

// --- Ports exposed to other domains ---

// HasActForTask satisfies task.ActStore: it gates marking a task completed.
func (s *Service) HasActForTask(ctx context.Context, taskID int64, taskType string) (bool, error) {
	switch taskType {
	case taskTypeInspection:
		return s.repo.ExistsInspectionByTaskID(ctx, taskID)
	case taskTypeReplacement:
		return s.repo.ExistsReplacementByTaskID(ctx, taskID)
	default:
		return false, nil
	}
}

// EnsureInspectionAct satisfies meter.ActStore and photo.ActStore: it lets
// those domains confirm tenant ownership of an inspection act before
// mutating its meters or photos.
func (s *Service) EnsureInspectionAct(ctx context.Context, actID, tenantID int64) error {
	_, err := s.repo.GetInspectionByID(ctx, actID, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrActNotFound
		}
		return err
	}
	return nil
}

// EnsureReplacementAct satisfies photo.ActStore.
func (s *Service) EnsureReplacementAct(ctx context.Context, actID, tenantID int64) error {
	_, err := s.repo.GetReplacementByID(ctx, actID, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrActNotFound
		}
		return err
	}
	return nil
}

// --- helpers ---

func (s *Service) editableTask(ctx context.Context, taskID, tenantID, userID int64, wantType string) (TaskInfo, error) {
	t, err := s.tasks.GetTask(ctx, taskID, tenantID)
	if err != nil {
		return TaskInfo{}, ErrTaskNotFound
	}
	if t.Type != wantType {
		return TaskInfo{}, ErrWrongTaskType
	}
	if t.AssigneeID == nil || *t.AssigneeID != userID {
		return TaskInfo{}, ErrNotAssignee
	}
	if t.Status != taskStatusInProgress {
		return TaskInfo{}, ErrTaskNotInProgress
	}
	return t, nil
}

func (s *Service) checkConsumer(ctx context.Context, tenantID int64, consumerID *int64) error {
	if consumerID == nil {
		return nil
	}
	if _, err := s.consumers.Get(ctx, *consumerID, tenantID); err != nil {
		return ErrConsumerNotFound
	}
	return nil
}

func (s *Service) enrichInspection(ctx context.Context, a InspectionAct) InspectionDTO {
	dto := toInspectionDTO(a)
	if addr, err := s.addresses.Get(ctx, a.AddressID, a.TenantID); err == nil {
		dto.AddressLabel = addr.Label
	}
	if a.ConsumerID != nil {
		if c, err := s.consumers.Get(ctx, *a.ConsumerID, a.TenantID); err == nil {
			dto.ConsumerName = c.Name
		}
	}
	if meters, err := s.meters.ListByAct(ctx, a.ID, a.TenantID); err == nil {
		dto.MeterCount = len(meters)
	}
	if photos, err := s.photos.ListForInspectionAct(ctx, a.ID); err == nil {
		dto.PhotoCount = len(photos)
	}
	return dto
}

func (s *Service) enrichReplacement(ctx context.Context, a ReplacementAct) ReplacementDTO {
	dto := toReplacementDTO(a)
	if addr, err := s.addresses.Get(ctx, a.AddressID, a.TenantID); err == nil {
		dto.AddressLabel = addr.Label
	}
	if photos, err := s.photos.ListForReplacementAct(ctx, a.ID); err == nil {
		dto.PhotoCount = len(photos)
	}
	return dto
}
