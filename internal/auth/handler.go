package auth

import (
	"encoding/json"
	"errors"
	"net/http"

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
	r.Post("/register-company", h.registerCompany)
	r.Post("/register-client", h.registerClient)
	r.Post("/register-master", h.registerMaster)
	r.Post("/bootstrap-super-admin", h.bootstrapSuperAdmin)
	r.Post("/login", h.login)
	r.Post("/refresh", h.refresh)
	r.Post("/logout", h.logout)
	return r
}

// UserRoutes exposes tenant staff management (dispatchers, electricians).
// Mounted separately at /api/users since, unlike Routes(), every endpoint
// here requires an authenticated tenant user.
func (h *Handler) UserRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.listUsers)
	r.With(middleware.RequireRoles(string(RoleTenantAdmin))).Post("/", h.createUser)
	return r
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := middleware.TenantIDFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", "Tenant context required")
		return
	}
	page, pageSize := httpx.ParsePageParams(r)
	role := Role(r.URL.Query().Get("role"))

	users, err := h.service.ListTenantUsers(r.Context(), tenantID, role, page, pageSize)
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list users")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, users)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := middleware.TenantIDFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", "Tenant context required")
		return
	}
	var req CreateTenantUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	user, err := h.service.CreateTenantUser(r.Context(), tenantID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, user)
}

func (h *Handler) registerCompany(w http.ResponseWriter, r *http.Request) {
	var req CompanyRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}

	resp, err := h.service.RegisterCompany(r.Context(), req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) registerClient(w http.ResponseWriter, r *http.Request) {
	var req RegisterMarketplaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}

	resp, err := h.service.RegisterClient(r.Context(), req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) registerMaster(w http.ResponseWriter, r *http.Request) {
	var req RegisterMarketplaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}

	resp, err := h.service.RegisterMaster(r.Context(), req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) bootstrapSuperAdmin(w http.ResponseWriter, r *http.Request) {
	var req BootstrapSuperAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}

	resp, err := h.service.BootstrapSuperAdmin(r.Context(), req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}

	resp, err := h.service.Login(r.Context(), req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}

	resp, err := h.service.Refresh(r.Context(), req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}

	if err := h.service.Logout(r.Context(), req); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", err.Error())
	case errors.Is(err, ErrInvalidRefreshToken):
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", err.Error())
	case errors.Is(err, ErrSuperAdminExists),
		errors.Is(err, ErrUsernameInUse),
		errors.Is(err, ErrTenantNameInUse),
		errors.Is(err, ErrTenantCodeInUse),
		errors.Is(err, ErrUserLimitExceeded):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	default:
		var validationErr interface{ Error() string }
		if errors.As(err, &validationErr) && err.Error() != "" && !errors.Is(err, ErrInvalidCredentials) {
			httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
			return
		}
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Authentication failed")
	}
}
