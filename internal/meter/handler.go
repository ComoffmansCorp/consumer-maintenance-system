package meter

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/middleware"
)

// Handler is mounted by the composition root under
// /api/acts/inspection/{id}/meters (the same {id} segment act's own
// inspection routes use, so chi's router tree stays consistent), so every
// route here relies on that "id" URL param already resolved by the parent.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Patch("/{meterID}", h.update)
	r.Delete("/{meterID}", h.delete)
	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tenantID, actID, ok := h.pathContext(w, r)
	if !ok {
		return
	}
	meters, err := h.service.ListByAct(r.Context(), actID, tenantID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, meters)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	tenantID, actID, ok := h.pathContext(w, r)
	if !ok {
		return
	}
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.Create(r.Context(), actID, tenantID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	tenantID, actID, ok := h.pathContext(w, r)
	if !ok {
		return
	}
	meterID, err := strconv.ParseInt(chi.URLParam(r, "meterID"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid meter id")
		return
	}
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.Update(r.Context(), meterID, actID, tenantID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	tenantID, actID, ok := h.pathContext(w, r)
	if !ok {
		return
	}
	meterID, err := strconv.ParseInt(chi.URLParam(r, "meterID"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid meter id")
		return
	}
	if err := h.service.Delete(r.Context(), meterID, actID, tenantID); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) pathContext(w http.ResponseWriter, r *http.Request) (tenantID, actID int64, ok bool) {
	tenantID, ok = middleware.TenantIDFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", "Tenant context required")
		return 0, 0, false
	}
	actID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid act id")
		return 0, 0, false
	}
	return tenantID, actID, true
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrMeterNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrInvalidType), errors.Is(err, ErrInvalidSealState), errors.Is(err, ErrSerialNumberRequired):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
	}
}
