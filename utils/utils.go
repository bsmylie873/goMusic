package utils

import (
	"context"
	"encoding/json"
	"goMusic/db"
	"goMusic/validation"
	"net/http"
	"time"
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

// ExecuteWithTransaction executes a SQL query with parameters inside a transaction with timeout
func ExecuteWithTransaction(w http.ResponseWriter, query string, args ...interface{}) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return false
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return false
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "Commit error: "+err.Error(), http.StatusInternalServerError)
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
