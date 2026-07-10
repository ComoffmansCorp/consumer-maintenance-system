package organization

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
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
	r.Post("/", h.create)
	r.Patch("/{tenantID}/plan", h.updatePlan)
	r.Delete("/{tenantID}", h.deactivate)
	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.service.List(r.Context())
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list tenants")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, tenants)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	if req.Plan == "" {
		req.Plan = PlanFree
	}

	tenant, err := h.service.Create(r.Context(), req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, tenant)
}

func (h *Handler) updatePlan(w http.ResponseWriter, r *http.Request) {
	tenantID, err := strconv.ParseInt(chi.URLParam(r, "tenantID"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid tenant id")
		return
	}

	var req UpdateTenantPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}

	tenant, err := h.service.UpdatePlan(r.Context(), tenantID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, tenant)
}

func (h *Handler) deactivate(w http.ResponseWriter, r *http.Request) {
	tenantID, err := strconv.ParseInt(chi.URLParam(r, "tenantID"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid tenant id")
		return
	}

	tenant, err := h.service.Deactivate(r.Context(), tenantID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, tenant)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrTenantNameInUse), errors.Is(err, ErrTenantCodeInUse), errors.Is(err, ErrInvalidTenantPlan):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
	}
}
