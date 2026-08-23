package httpadapter

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e APIError) Error() string { return e.Message }
func classify(err error) APIError {
	if err == nil {
		return APIError{Code: "ok", Status: 200}
	}
	if errors.Is(err, contextCanceled{}) {
		return APIError{Code: "cancelled", Message: err.Error(), Status: 499}
	}
	return APIError{Code: "invalid_request", Message: err.Error(), Status: 400}
}

type contextCanceled struct{}

func (contextCanceled) Error() string { return "context canceled" }
func writeError(w http.ResponseWriter, err error) {
	api := classify(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(api.Status)
	_ = json.NewEncoder(w).Encode(api)
}
func isJSON(r *http.Request) bool { return strings.Contains(r.Header.Get("Content-Type"), "json") }
