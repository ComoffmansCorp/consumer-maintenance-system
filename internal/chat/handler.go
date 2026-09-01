package chat

import (
	"encoding/json"
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
	r.Get("/requests/{id}/messages", h.listMessages)
	r.Post("/requests/{id}/messages", h.sendMessage)
	return r
}

func (h *Handler) listMessages(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	requestID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	var sinceID int64
	if raw := r.URL.Query().Get("sinceId"); raw != "" {
		sinceID, _ = strconv.ParseInt(raw, 10, 64)
	}
	messages, err := h.service.ListMessages(r.Context(), requestID, auth.UserID, sinceID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, messages)
}

func (h *Handler) sendMessage(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	requestID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.SendMessage(r.Context(), requestID, auth.UserID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrRequestNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrTextRequired):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	case errors.Is(err, ErrForbidden), errors.Is(err, ErrNotAssignedYet):
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", err.Error())
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
	}
}
