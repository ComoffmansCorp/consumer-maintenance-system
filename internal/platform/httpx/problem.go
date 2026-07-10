package httpx

import (
	"encoding/json"
	"net/http"
)

type Problem struct {
	Type     string         `json:"type,omitempty"`
	Title    string         `json:"title"`
	Status   int            `json:"status"`
	Detail   string         `json:"detail,omitempty"`
	Instance string         `json:"instance,omitempty"`
	Errors   map[string]any `json:"errors,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteProblem(w http.ResponseWriter, status int, title, detail string) {
	WriteJSON(w, status, Problem{
		Type:   "about:blank",
		Title:  title,
		Status: status,
		Detail: detail,
	})
}

func WriteValidationProblem(w http.ResponseWriter, detail string, errors map[string]any) {
	WriteJSON(w, http.StatusBadRequest, Problem{
		Type:   "about:blank",
		Title:  "Validation failed",
		Status: http.StatusBadRequest,
		Detail: detail,
		Errors: errors,
	})
}
