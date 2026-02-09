package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"setlist/api/middleware"
	"setlist/api/model"
	"setlist/api/repository"
	"setlist/api/service"
)

type UserHandler struct {
	UserService service.UserService
}

func (h UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var payload service.AuthPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.UserService.Signup(r.Context(), payload)
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
	json.NewEncoder(w).Encode(response)
}

func (h UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload service.LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.UserService.Login(r.Context(), payload)
	if err != nil {
		writeError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		writeError(w, "Could not identify user from token", http.StatusInternalServerError)
		return
	}

	var payload service.UpdatePasswordPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.NewPassword == "" {
		writeError(w, "New password cannot be empty", http.StatusBadRequest)
		return
	}

	err := h.UserService.UpdatePassword(r.Context(), userID, payload)
	if err != nil {
		if err.Error() == "invalid current password" {
			writeError(w, err.Error(), http.StatusUnauthorized)
			return
		}
		writeError(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}

func (h UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	users, err := h.UserService.SearchUsers(r.Context(), query)
	if err != nil {
		writeError(w, "Failed to search users", http.StatusInternalServerError)
		return
	}
	if users == nil {
		users = make([]model.User, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
