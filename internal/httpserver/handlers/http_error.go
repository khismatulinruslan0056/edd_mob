package handlers

import (
	"Effective_Mobile/internal/httpserver/handlers/dto"
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: msg})
}
