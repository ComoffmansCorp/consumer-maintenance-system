package catalog

import (
	"encoding/json"
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
	// Public (see router.go publicPaths) -- JWTAuth matches by exact path
	// only, not method, so admin writes live under the separate, shared
	// /api/admin namespace (AdminRoutes below) rather than on this same
	// "/categories"/"/services" path, or they would accidentally become
	// unauthenticated too.
	r.Get("/categories", h.listCategories)
	r.Get("/services", h.listServices)
	return r
}

// AdminRoutes is mounted separately at /api/admin (categories and services
// share this handler, so both resource prefixes are registered here).
func (h *Handler) AdminRoutes() chi.Router {
	admin := middleware.RequireRoles(roleSuperAdmin)

	r := chi.NewRouter()
	r.With(admin).Post("/categories", h.createCategory)
	r.With(admin).Put("/categories/{id}", h.updateCategory)
	r.With(admin).Post("/services", h.createService)
	r.With(admin).Put("/services/{id}", h.updateService)
	return r
}

func (h *Handler) listCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.ListCategories(r.Context())
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list categories")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, categories)
}

func (h *Handler) createCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.CreateCategory(r.Context(), req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto)
}

func (h *Handler) updateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid category id")
		return
	}
	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.UpdateCategory(r.Context(), id, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) listServices(w http.ResponseWriter, r *http.Request) {
	var categoryID int64
	if raw := r.URL.Query().Get("categoryId"); raw != "" {
		categoryID, _ = strconv.ParseInt(raw, 10, 64)
	}
	services, err := h.service.ListServices(r.Context(), categoryID)
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list services")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, services)
}

func (h *Handler) createService(w http.ResponseWriter, r *http.Request) {
	var req CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.CreateService(r.Context(), req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto)
}

func (h *Handler) updateService(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid service id")
		return
	}
	var req UpdateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.UpdateService(r.Context(), id, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrCategoryNotFound), errors.Is(err, ErrServiceNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrNameRequired):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
	}
}
