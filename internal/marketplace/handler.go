package marketplace

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
	roleMaster     = "MASTER"
	roleClient     = "CLIENT"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	admin := middleware.RequireRoles(roleSuperAdmin)
	master := middleware.RequireRoles(roleMaster)
	client := middleware.RequireRoles(roleClient)
	clientOrMaster := middleware.RequireRoles(roleClient, roleMaster)

	r := chi.NewRouter()
	// GET here is public (see router.go publicPaths) -- JWTAuth matches by
	// exact path only, not method, so admin writes MUST live on a different
	// path (/admin/...) rather than POST on the same "/categories"/"/services"
	// path, or they would accidentally become unauthenticated too.
	r.Get("/categories", h.listCategories)
	r.Get("/services", h.listServices)
	r.With(admin).Post("/admin/categories", h.createCategory)
	r.With(admin).Post("/admin/services", h.createService)

	r.With(master).Get("/master/profile", h.getMasterProfile)
	r.With(master).Put("/master/profile", h.updateMasterProfile)

	r.With(client).Post("/requests", h.createRequest)
	r.With(clientOrMaster).Get("/requests", h.listMyRequests)
	r.With(master).Get("/requests/open", h.listOpenRequests)
	r.With(clientOrMaster).Get("/requests/{id}", h.getRequest)
	r.With(master).Post("/requests/{id}/claim", h.claim)
	r.With(master).Post("/requests/{id}/complete", h.complete)
	r.With(clientOrMaster).Post("/requests/{id}/cancel", h.cancel)
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
	dto, err := h.service.CreateCategory(r.Context(), req.Name)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto)
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
	var req CreateOfferingRequest
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

func (h *Handler) getMasterProfile(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	dto, err := h.service.GetMasterProfile(r.Context(), auth.UserID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) updateMasterProfile(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	var req UpdateMasterProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.UpdateMasterProfile(r.Context(), auth.UserID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) createRequest(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	var req CreateRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.CreateRequest(r.Context(), auth.UserID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto)
}

// listMyRequests branches by role: a client sees requests they created, a
// master sees requests currently or previously assigned to them. The two
// never overlap for a single account since a user holds exactly one role.
func (h *Handler) listMyRequests(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	page, pageSize := httpx.ParsePageParams(r)

	var result httpx.Page[RequestDTO]
	var err error
	if auth.Role == roleMaster {
		result, err = h.service.ListClaimedByMaster(r.Context(), auth.UserID, page, pageSize)
	} else {
		result, err = h.service.ListMine(r.Context(), auth.UserID, page, pageSize)
	}
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list requests")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) listOpenRequests(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	page, pageSize := httpx.ParsePageParams(r)
	result, err := h.service.ListOpenForMaster(r.Context(), auth.UserID, page, pageSize)
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list open requests")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) getRequest(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	id, err := parseID(r)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	dto, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	if dto.ClientID != auth.UserID && (dto.MasterID == nil || *dto.MasterID != auth.UserID) {
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", "Not your request")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) claim(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	id, err := parseID(r)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	dto, err := h.service.Claim(r.Context(), id, auth.UserID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) complete(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	id, err := parseID(r)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	dto, err := h.service.Complete(r.Context(), id, auth.UserID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	id, err := parseID(r)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	var req CancelRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.Cancel(r.Context(), id, auth.UserID, req.Reason)
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
	case errors.Is(err, ErrRequestNotFound), errors.Is(err, ErrCategoryNotFound), errors.Is(err, ErrServiceNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrDescriptionRequired), errors.Is(err, ErrAddressRequired), errors.Is(err, ErrCancelReasonRequired):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	case errors.Is(err, ErrNotSpecialized):
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", err.Error())
	case errors.Is(err, ErrInvalidTransition):
		httpx.WriteProblem(w, http.StatusConflict, "Conflict", err.Error())
	case errors.Is(err, ErrNotOwner), errors.Is(err, ErrNotAssignedMaster):
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", err.Error())
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
	}
}
