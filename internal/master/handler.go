package master

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/middleware"
)

const (
	roleSuperAdmin = "SUPER_ADMIN"
	roleMaster     = "MASTER"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	master := middleware.RequireRoles(roleMaster)

	r := chi.NewRouter()
	r.With(master).Get("/profile", h.getProfile)
	r.With(master).Put("/profile", h.updateProfile)
	return r
}

// AdminRoutes is mounted separately at /api/admin/masters.
func (h *Handler) AdminRoutes() chi.Router {
	admin := middleware.RequireRoles(roleSuperAdmin)

	r := chi.NewRouter()
	r.With(admin).Get("/", h.listProfiles)
	return r
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	dto, err := h.service.GetProfile(r.Context(), auth.UserID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.UpdateProfile(r.Context(), auth.UserID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) listProfiles(w http.ResponseWriter, r *http.Request) {
	page, pageSize := httpx.ParsePageParams(r)
	result, err := h.service.ListProfiles(r.Context(), page, pageSize)
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list masters")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrProfileNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrInvalidSpecialization):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
	}
}
