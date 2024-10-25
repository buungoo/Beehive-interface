package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code,omitempty"`
}

// This is used for every successfull request
func SendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		SendErrorResponse(w, "Could not encode JSON response", http.StatusInternalServerError)
	}
}

// This is used for every error response
func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Genereate error message in JSON format
	errorResponse := ErrorResponse{
		Error: message,
		Code:  statusCode,
	}

	json.NewEncoder(w).Encode(errorResponse)

}
