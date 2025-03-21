package utils

import (
	"encoding/json"
	"goMusic/validation"
	"net/http"
)

// DecodeJSONBody decodes a JSON request body into the provided struct
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid request"})
		return false
	}
	return true
}

// ValidateRequestBody validates the provided struct and handles error responses
func ValidateRequestBody(w http.ResponseWriter, v interface{}) bool {
	if err := validation.ValidateStruct(v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return false
	}
	return true
}

// DecodeAndValidate combines both operations for common use case
func DecodeAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	return DecodeJSONBody(w, r, v) && ValidateRequestBody(w, v)
}
