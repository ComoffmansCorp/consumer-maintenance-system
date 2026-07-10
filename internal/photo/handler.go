package photo

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/middleware"
)

const maxUploadMemory = 32 << 20 // buffer in memory before spilling to temp files

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.upload)
	r.Get("/inspection/{actID}", h.listInspection)
	r.Get("/replacement/{actID}", h.listReplacement)
	r.Get("/{id}", h.download)
	r.Delete("/{id}", h.delete)
	return r
}

func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.authContext(w, r)
	if !ok {
		return
	}
	if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid multipart form")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "File is required")
		return
	}
	defer func() { _ = file.Close() }()

	in := UploadInput{
		TenantID:         tenantID,
		UploadedBy:       userID,
		Note:             r.FormValue("note"),
		OriginalFilename: header.Filename,
		ContentType:      header.Header.Get("Content-Type"),
		Size:             header.Size,
		Content:          file,
	}
	if raw := r.FormValue("inspectionActId"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid inspectionActId")
			return
		}
		in.InspectionActID = &id
	}
	if raw := r.FormValue("replacementActId"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid replacementActId")
			return
		}
		in.ReplacementActID = &id
	}

	dto, err := h.service.Upload(r.Context(), in)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto)
}

func (h *Handler) listInspection(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	actID, err := strconv.ParseInt(chi.URLParam(r, "actID"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid act id")
		return
	}
	photos, err := h.service.ListInspection(r.Context(), actID, tenantID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, photos)
}

func (h *Handler) listReplacement(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	actID, err := strconv.ParseInt(chi.URLParam(r, "actID"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid act id")
		return
	}
	photos, err := h.service.ListReplacement(r.Context(), actID, tenantID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, photos)
}

func (h *Handler) download(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid photo id")
		return
	}
	p, path, err := h.service.GetFile(r.Context(), id, tenantID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", p.ContentType)
	http.ServeFile(w, r, path)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	tenantID, _, ok := h.authContext(w, r)
	if !ok {
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", "Invalid photo id")
		return
	}
	if err := h.service.Delete(r.Context(), id, tenantID); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
	case errors.Is(err, ErrPhotoNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrInvalidTarget), errors.Is(err, ErrFileRequired), errors.Is(err, ErrFileTooLarge),
		errors.Is(err, ErrUnsupportedType):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad request", err.Error())
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal server error", "Operation failed")
	}
}
