package request

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
	master := middleware.RequireRoles(roleMaster)
	client := middleware.RequireRoles(roleClient)
	clientOrMaster := middleware.RequireRoles(roleClient, roleMaster)

	r := chi.NewRouter()
	r.With(client).Post("/", h.createRequest)
	r.With(clientOrMaster).Get("/", h.listMyRequests)
	r.With(master).Get("/open", h.listOpenRequests)
	r.With(client).Get("/favorites", h.listFavorites)
	r.With(client).Post("/favorites", h.addFavorite)
	r.With(client).Delete("/favorites/{masterId}", h.removeFavorite)
	r.With(clientOrMaster).Get("/{id}", h.getRequest)
	r.With(master).Post("/{id}/offers", h.submitOffer)
	r.With(client).Get("/{id}/offers", h.listOffers)
	r.With(client).Post("/{id}/offers/{offerId}/accept", h.acceptOffer)
	r.With(master).Post("/{id}/complete", h.complete)
	r.With(clientOrMaster).Post("/{id}/cancel", h.cancel)
	return r
}

// AdminRoutes is mounted separately at /api/admin/requests.
func (h *Handler) AdminRoutes() chi.Router {
	admin := middleware.RequireRoles(roleSuperAdmin)

	r := chi.NewRouter()
	r.With(admin).Get("/", h.listAdmin)
	return r
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
// master sees requests currently or previously assigned to them.
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

func (h *Handler) listAdmin(w http.ResponseWriter, r *http.Request) {
	page, pageSize := httpx.ParsePageParams(r)
	status := r.URL.Query().Get("status")
	result, err := h.service.ListAdmin(r.Context(), status, page, pageSize)
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list requests")
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
	id, err := parseID(r, "id")
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	dto, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	isParticipant := dto.ClientID == auth.UserID || (dto.MasterID != nil && *dto.MasterID == auth.UserID)
	// A master browsing the open pool hasn't been assigned yet (MasterID is
	// nil) -- they still need to see the request's details to decide
	// whether to submit an offer, so OPEN requests are readable by any
	// master, not just participants. Once assigned/completed/canceled, only
	// the client and the assigned master can see it.
	isMasterBrowsingOpen := auth.Role == roleMaster && dto.Status == StatusOpen
	if !isParticipant && !isMasterBrowsingOpen {
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", "Not your request")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) submitOffer(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	id, err := parseID(r, "id")
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	var req SubmitOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.SubmitOffer(r.Context(), id, auth.UserID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto)
}

func (h *Handler) listOffers(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	id, err := parseID(r, "id")
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	offers, err := h.service.ListOffers(r.Context(), id, auth.UserID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, offers)
}

func (h *Handler) acceptOffer(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	id, err := parseID(r, "id")
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	offerID, err := parseID(r, "offerId")
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid offer id")
		return
	}
	dto, err := h.service.AcceptOffer(r.Context(), id, offerID, auth.UserID)
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
	id, err := parseID(r, "id")
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	dto, err := h.service.CompleteRequest(r.Context(), id, auth.UserID)
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
	id, err := parseID(r, "id")
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid request id")
		return
	}
	var req CancelRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	dto, err := h.service.CancelRequest(r.Context(), id, auth.UserID, req.Reason)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

func (h *Handler) listFavorites(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	favorites, err := h.service.ListFavorites(r.Context(), auth.UserID)
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to list favorites")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, favorites)
}

type addFavoriteRequest struct {
	MasterID int64 `json:"masterId"`
}

func (h *Handler) addFavorite(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	var req addFavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid JSON body")
		return
	}
	if err := h.service.AddFavorite(r.Context(), auth.UserID, req.MasterID); err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to add favorite")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) removeFavorite(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
		return
	}
	masterID, err := parseID(r, "masterId")
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid master id")
		return
	}
	if err := h.service.RemoveFavorite(r.Context(), auth.UserID, masterID); err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Failed to remove favorite")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(r *http.Request, param string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, param), 10, 64)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrRequestNotFound), errors.Is(err, ErrServiceNotFound), errors.Is(err, ErrOfferNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrDescriptionRequired), errors.Is(err, ErrAddressRequired), errors.Is(err, ErrCancelReasonRequired),
		errors.Is(err, ErrPriceInvalid), errors.Is(err, ErrOfferAlreadyExists):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	case errors.Is(err, ErrNotSpecialized), errors.Is(err, ErrNotOwner), errors.Is(err, ErrNotAssignedMaster):
		httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", err.Error())
	case errors.Is(err, ErrInvalidTransition), errors.Is(err, ErrOfferNotPending):
		httpx.WriteProblem(w, http.StatusConflict, "Conflict", err.Error())
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
	}
}
