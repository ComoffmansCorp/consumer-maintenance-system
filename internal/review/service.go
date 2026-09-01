package review

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

var (
	ErrNotEligible         = errors.New("request is not eligible for a review")
	ErrRatingInvalid       = errors.New("rating must be between 1 and 5")
	ErrReviewAlreadyExists = errors.New("a review already exists for this request")
	ErrReviewNotFound      = errors.New("review not found")
)

type Service struct {
	repo     *Repository
	requests RequestPort
	masters  MasterPort
}

func NewService(repo *Repository, requests RequestPort, masters MasterPort) *Service {
	return &Service{repo: repo, requests: requests, masters: masters}
}

// CreateReview requires the request be COMPLETED and owned by the reviewing
// client (checked via RequestPort), then folds the rating into the master's
// running average (via MasterPort) -- done as two sequential calls, not one
// cross-domain transaction: each domain owns its own atomicity, and the
// rating update is itself already atomic on the master side.
func (s *Service) CreateReview(ctx context.Context, clientID int64, req CreateReviewRequest) (ReviewDTO, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return ReviewDTO{}, ErrRatingInvalid
	}

	masterID, eligible, err := s.requests.GetCompletedForReview(ctx, req.RequestID, clientID)
	if err != nil {
		return ReviewDTO{}, err
	}
	if !eligible {
		return ReviewDTO{}, ErrNotEligible
	}

	created, err := s.repo.Create(ctx, req.RequestID, clientID, masterID, req.Rating, strings.TrimSpace(req.Comment))
	if err != nil {
		if isUniqueViolation(err) {
			return ReviewDTO{}, ErrReviewAlreadyExists
		}
		return ReviewDTO{}, fmt.Errorf("create review: %w", err)
	}

	if err := s.masters.RecordReview(ctx, masterID, req.Rating); err != nil {
		return ReviewDTO{}, fmt.Errorf("record review on master profile: %w", err)
	}

	return ToReviewDTO(created), nil
}

func (s *Service) ListVisibleByMaster(ctx context.Context, masterID int64, page, pageSize int) (httpx.Page[ReviewDTO], error) {
	limit, offset := int32(pageSize), httpx.Offset(page, pageSize)
	items, err := s.repo.ListVisibleByMaster(ctx, masterID, limit, offset)
	if err != nil {
		return httpx.Page[ReviewDTO]{}, err
	}
	total, err := s.repo.CountVisibleByMaster(ctx, masterID)
	if err != nil {
		return httpx.Page[ReviewDTO]{}, err
	}
	out := make([]ReviewDTO, 0, len(items))
	for _, r := range items {
		out = append(out, ToReviewDTO(r))
	}
	return httpx.NewPage(out, page, pageSize, total), nil
}

func (s *Service) HideReview(ctx context.Context, id int64) (ReviewDTO, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return ReviewDTO{}, ErrReviewNotFound
	}
	r, err := s.repo.Hide(ctx, id)
	if err != nil {
		return ReviewDTO{}, fmt.Errorf("hide review: %w", err)
	}
	return ToReviewDTO(r), nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
