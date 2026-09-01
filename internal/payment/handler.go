package payment

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/middleware"
)

const roleSuperAdmin = "SUPER_ADMIN"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/requests/{id}/payment", h.getForRequest)
	return r
}

// AdminRoutes is mounted separately at /api/admin/payments.
func (h *Handler) AdminRoutes() chi.Router {
	admin := middleware.RequireRoles(roleSuperAdmin)

	r := chi.NewRouter()
	r.With(admin).Get("/", h.listAdmin)
	return r
}

func (h *Handler) getForRequest(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	requestID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	dto, err := h.service.GetForRequest(r.Context(), requestID, auth.UserID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) listAdmin(w http.ResponseWriter, r *http.Request) {
	page, pageSize := httpx.ParsePageParams(r)
	result, err := h.service.ListAdmin(r.Context(), page, pageSize)
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list payments")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrPaymentNotFound), errors.Is(err, ErrRequestNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrForbidden):
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", err.Error())
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
	}
}
