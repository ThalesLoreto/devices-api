package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// WriteJSONResponse writes a JSON response with the given status code and data
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// WriteErrorResponse writes an error response with the given status code and message
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	errorResp := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	WriteJSONResponse(w, statusCode, errorResp)
}
