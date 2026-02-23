package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"setlist/api/apierror"
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
		writeAppError(w, apierror.InvalidRequest("Corps de la requête invalide."))
		return
	}

	response, err := h.UserService.Signup(r.Context(), payload)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateUsername) {
			writeAppError(w, apierror.UsernameTaken())
			return
		}
		if errors.Is(err, repository.ErrDuplicateBand) {
			writeAppError(w, apierror.BandNameTaken())
			return
		}
		var appErr *apierror.AppError
		if errors.As(err, &appErr) {
			writeAppError(w, appErr)
			return
		}
		writeAppError(w, apierror.InternalError("inscription"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload service.LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeAppError(w, apierror.InvalidRequest("Corps de la requête invalide."))
		return
	}

	response, err := h.UserService.Login(r.Context(), payload)
	if err != nil {
		writeAppError(w, apierror.InvalidCredentials())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		writeAppError(w, apierror.NewServerError(apierror.ErrInternal, "Impossible d'identifier l'utilisateur depuis le token."))
		return
	}

	var payload service.UpdatePasswordPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeAppError(w, apierror.InvalidRequest("Corps de la requête invalide."))
		return
	}

	if payload.NewPassword == "" {
		writeAppError(w, apierror.InvalidRequest("Le nouveau mot de passe ne peut pas être vide."))
		return
	}

	err := h.UserService.UpdatePassword(r.Context(), userID, payload)
	if err != nil {
		if err.Error() == "invalid current password" {
			writeAppError(w, apierror.WrongCurrentPassword())
			return
		}
		var appErr *apierror.AppError
		if errors.As(err, &appErr) {
			writeAppError(w, appErr)
			return
		}
		writeAppError(w, apierror.InternalError("mise à jour du mot de passe"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Mot de passe mis à jour avec succès."})
}

func (h UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	users, err := h.UserService.SearchUsers(r.Context(), query)
	if err != nil {
		writeAppError(w, apierror.InternalError("recherche d'utilisateurs"))
		return
	}
	if users == nil {
		users = make([]model.User, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
