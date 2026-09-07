package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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
	r.Post("/register-client", h.registerClient)
	r.Post("/register-master", h.registerMaster)
	r.Post("/bootstrap-super-admin", h.bootstrapSuperAdmin)
	r.Post("/login", h.login)
	r.Post("/refresh", h.refresh)
	r.Post("/logout", h.logout)
	return r
}

func (h *Handler) registerClient(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
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
	var req RegisterRequest
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

	if err := h.service.Logout(r.Context(), req, bearerToken(r)); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// bearerToken extracts the access token from the Authorization header, if
// present. /api/auth/logout is a public route (see router.go's
// publicPaths), so it never runs through middleware.JWTAuth -- the client
// may or may not still send a bearer token here, and either is valid.
func bearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrInvalidRefreshToken):
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", err.Error())
	case errors.Is(err, ErrSuperAdminExists), errors.Is(err, ErrUsernameInUse):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	default:
		var validationErr interface{ Error() string }
		if errors.As(err, &validationErr) && err.Error() != "" {
			httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
			return
		}
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Authentication failed")
	}
}
