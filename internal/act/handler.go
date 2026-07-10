package act

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/middleware"
)

const roleElectrician = "ELECTRICIAN"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	inspect := middleware.RequireRoles(roleElectrician)

	r := chi.NewRouter()
	r.Route("/inspection", func(ir chi.Router) {
		ir.With(inspect).Post("/", h.createInspection)
		ir.Get("/by-task/{taskID}", h.getInspectionByTask)
		ir.Get("/{id}", h.getInspection)
		ir.With(inspect).Patch("/{id}", h.updateInspection)
		ir.Get("/{id}/pdf", h.inspectionPDF)
	})
	r.Route("/replacement", func(rr chi.Router) {
		rr.With(inspect).Post("/", h.createReplacement)
		rr.Get("/by-task/{taskID}", h.getReplacementByTask)
		rr.Get("/{id}", h.getReplacement)
		rr.With(inspect).Patch("/{id}", h.updateReplacement)
		rr.Get("/{id}/pdf", h.replacementPDF)
	})
	return r
}

func (h *Handler) createInspection(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.authContext(w, r)
	if !ok {
		return
	}
	var req CreateInspectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.CreateInspection(r.Context(), tenantID, userID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto)
}

func (h *Handler) getInspection(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid act id")
		return
	}
	dto, err := h.service.GetInspection(r.Context(), id, tenantID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) getInspectionByTask(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid task id")
		return
	}
	dto, err := h.service.GetInspectionByTask(r.Context(), taskID, tenantID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) updateInspection(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid act id")
		return
	}
	var req UpdateInspectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.UpdateInspection(r.Context(), id, tenantID, userID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) inspectionPDF(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid act id")
		return
	}
	data, filename, err := h.service.GenerateInspectionPDF(r.Context(), id, tenantID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writePDF(w, data, filename)
}

func (h *Handler) createReplacement(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.authContext(w, r)
	if !ok {
		return
	}
	var req CreateReplacementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.CreateReplacement(r.Context(), tenantID, userID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto)
}

func (h *Handler) getReplacement(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid act id")
		return
	}
	dto, err := h.service.GetReplacement(r.Context(), id, tenantID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) getReplacementByTask(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid task id")
		return
	}
	dto, err := h.service.GetReplacementByTask(r.Context(), taskID, tenantID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) updateReplacement(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid act id")
		return
	}
	var req UpdateReplacementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.UpdateReplacement(r.Context(), id, tenantID, userID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) replacementPDF(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid act id")
		return
	}
	data, filename, err := h.service.GenerateReplacementPDF(r.Context(), id, tenantID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writePDF(w, data, filename)
}

func writePDF(w http.ResponseWriter, data []byte, filename string) {
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (h *Handler) authContext(w http.ResponseWriter, r *http.Request) (tenantID, userID int64, ok bool) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok || auth.TenantID == nil {
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", "Tenant context required")
		return 0, 0, false
	}
	return *auth.TenantID, auth.UserID, true
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrActNotFound), errors.Is(err, ErrTaskNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrWrongTaskType), errors.Is(err, ErrConsumerNotFound), errors.Is(err, ErrInvalidInspectionType),
		errors.Is(err, ErrAccountNumberRequired):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	case errors.Is(err, ErrActAlreadyExists), errors.Is(err, ErrTaskNotInProgress):
		httpx.WriteProblem(w, http.StatusConflict, "Conflict", err.Error())
	case errors.Is(err, ErrNotAssignee):
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", err.Error())
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
	}
}
