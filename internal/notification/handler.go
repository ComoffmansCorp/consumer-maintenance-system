package notification

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Get("/unread-count", h.unreadCount)
	r.Post("/{id}/read", h.markRead)
	r.Post("/read-all", h.markAllRead)
	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.authContext(w, r)
	if !ok {
		return
	}
	page, pageSize := httpx.ParsePageParams(r)
	unreadOnly := r.URL.Query().Get("unread") == "true"

	result, err := h.service.List(r.Context(), tenantID, userID, unreadOnly, page, pageSize)
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list notifications")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) unreadCount(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.authContext(w, r)
	if !ok {
		return
	}
	count, err := h.service.UnreadCount(r.Context(), tenantID, userID)
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to count notifications")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]int64{"unreadCount": count})
}

func (h *Handler) markRead(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid notification id")
		return
	}
	dto, err := h.service.MarkRead(r.Context(), id, tenantID, userID)
	if err != nil {
		if errors.Is(err, ErrNotificationNotFound) {
			httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
			return
		}
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) markAllRead(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.authContext(w, r)
	if !ok {
		return
	}
	if err := h.service.MarkAllRead(r.Context(), tenantID, userID); err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) authContext(w http.ResponseWriter, r *http.Request) (tenantID, userID int64, ok bool) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok || auth.TenantID == nil {
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", "Tenant context required")
		return 0, 0, false
	}
	return *auth.TenantID, auth.UserID, true
}
