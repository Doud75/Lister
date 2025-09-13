package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"setlist/api/repository"
	"setlist/api/service"
)

type UserHandler struct {
	UserService service.UserService
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var payload service.AuthPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.UserService.Signup(r.Context(), payload)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateBand) || errors.Is(err, repository.ErrDuplicateUsername) {
			writeError(w, err.Error(), http.StatusConflict)
			return
		}
		writeError(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h UserHandler) Join(w http.ResponseWriter, r *http.Request) {
	var payload service.AuthPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.UserService.Join(r.Context(), payload)
	if err != nil {
		if err.Error() == "band not found" {
			writeError(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrDuplicateUsername) {
			writeError(w, err.Error(), http.StatusConflict)
			return
		}
		writeError(w, "Failed to join band", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload service.LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.UserService.Login(r.Context(), payload)
	if err != nil {
		writeError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
