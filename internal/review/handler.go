package review

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/middleware"
)

const (
	roleSuperAdmin = "SUPER_ADMIN"
	roleClient     = "CLIENT"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Routes is mounted at /api/reviews, behind JWTAuth like every other
// protected domain route.
func (h *Handler) Routes() chi.Router {
	client := middleware.RequireRoles(roleClient)

	r := chi.NewRouter()
	r.With(client).Post("/", h.createReview)
	return r
}

// AdminRoutes is mounted separately at /api/admin/reviews.
func (h *Handler) AdminRoutes() chi.Router {
	admin := middleware.RequireRoles(roleSuperAdmin)

	r := chi.NewRouter()
	r.With(admin).Put("/{id}/hide", h.hideReview)
	return r
}

// ListMasterReviews backs the public GET /api/masters/{id}/reviews route.
// It is wired directly on the root router in server/router.go, bypassing
// JWTAuth entirely -- the path carries a {id} segment, and JWTAuth's
// publicPaths map only matches exact literal paths, so this route cannot be
// listed there the way the static catalog GETs are.
func (h *Handler) ListMasterReviews(w http.ResponseWriter, r *http.Request) {
	masterID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid master id")
		return
	}
	page, pageSize := httpx.ParsePageParams(r)
	result, err := h.service.ListVisibleByMaster(r.Context(), masterID, page, pageSize)
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list reviews")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) createReview(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	var req CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.CreateReview(r.Context(), auth.UserID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto)
}

func (h *Handler) hideReview(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid review id")
		return
	}
	dto, err := h.service.HideReview(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrReviewNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrRatingInvalid), errors.Is(err, ErrReviewAlreadyExists):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	case errors.Is(err, ErrNotEligible):
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", err.Error())
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
	}
}
