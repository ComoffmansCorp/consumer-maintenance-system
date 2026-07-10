package task

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
	roleDispatcher  = "DISPATCHER"
	roleTenantAdmin = "TENANT_ADMIN"
	roleElectrician = "ELECTRICIAN"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	dispatch := middleware.RequireRoles(roleDispatcher, roleTenantAdmin)
	inspect := middleware.RequireRoles(roleElectrician)

	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Get("/{id}", h.get)
	r.With(dispatch).Post("/", h.create)
	r.With(dispatch).Post("/{id}/assign", h.assign)
	r.With(dispatch).Post("/{id}/cancel", h.cancel)
	r.With(inspect).Post("/{id}/start", h.start)
	r.With(inspect).Post("/{id}/complete", h.complete)
	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tenantID, auth, ok := h.authContext(w, r)
	if !ok {
		return
	}
	page, pageSize := httpx.ParsePageParams(r)
	status := Status(r.URL.Query().Get("status"))
	typ := Type(r.URL.Query().Get("type"))

	var assigneeID int64
	if raw := r.URL.Query().Get("assigneeId"); raw != "" {
		assigneeID, _ = strconv.ParseInt(raw, 10, 64)
	}
	// Inspectors only ever see their own tasks, regardless of the query.
	if auth.Role == roleElectrician {
		assigneeID = auth.UserID
	}

	result, err := h.service.List(r.Context(), tenantID, status, typ, assigneeID, page, pageSize)
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list tasks")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.Create(r.Context(), tenantID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid task id")
		return
	}
	dto, err := h.service.Get(r.Context(), id, tenantID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) assign(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid task id")
		return
	}
	var req AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.Assign(r.Context(), id, tenantID, req.AssigneeID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) start(w http.ResponseWriter, r *http.Request) {
	tenantID, auth, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid task id")
		return
	}
	dto, err := h.service.Start(r.Context(), id, tenantID, auth.UserID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) complete(w http.ResponseWriter, r *http.Request) {
	tenantID, auth, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid task id")
		return
	}
	dto, err := h.service.Complete(r.Context(), id, tenantID, auth.UserID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid task id")
		return
	}
	var req StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.Cancel(r.Context(), id, tenantID, req.Reason)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) authContext(w http.ResponseWriter, r *http.Request) (tenantID int64, auth middleware.AuthContext, ok bool) {
	auth, ok = middleware.AuthFromContext(r.Context())
	if !ok || auth.TenantID == nil {
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", "Tenant context required")
		return 0, middleware.AuthContext{}, false
	}
	return *auth.TenantID, auth, true
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTaskNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrAddressNotFound), errors.Is(err, ErrAssigneeNotFound), errors.Is(err, ErrInvalidType),
		errors.Is(err, ErrCancelReasonRequired), errors.Is(err, ErrActNotFilled):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	case errors.Is(err, ErrInvalidTransition):
		httpx.WriteProblem(w, http.StatusConflict, "Conflict", err.Error())
	case errors.Is(err, ErrNotAssignee):
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", err.Error())
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
	}
}
