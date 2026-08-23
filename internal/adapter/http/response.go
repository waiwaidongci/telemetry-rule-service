package httpadapter

import (
	"encoding/json"
	"net/http"
	"time"
)

type EnvelopeResponse struct {
	Data      any       `json:"data"`
	RequestID string    `json:"request_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(EnvelopeResponse{Data: data, Timestamp: time.Now().UTC()})
}
func respondList(w http.ResponseWriter, data any, total int) {
	w.Header().Set("X-Total-Count", itoa(total))
	respond(w, http.StatusOK, data)
}
func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	digits := ""
	for value > 0 {
		digits = string(byte('0'+value%10)) + digits
		value /= 10
	}
	if negative {
		digits = "-" + digits
	}
	return digits
}
