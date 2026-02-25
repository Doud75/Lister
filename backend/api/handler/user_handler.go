package handler

import (
	"errors"
	"net/http"
	"setlist/api/apierror"
	"setlist/api/model"
	"setlist/api/repository"
	"setlist/api/service"
)

type UserHandler struct {
	UserService service.UserService
}

func (h UserHandler) Signup(w http.ResponseWriter, r *http.Request) error {
	payload, err := DecodeJSON[service.AuthPayload](r)
	if err != nil {
		return err
	}

	response, err := h.UserService.Signup(r.Context(), payload)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateUsername) {
			return apierror.UsernameTaken()
		}
		if errors.Is(err, repository.ErrDuplicateBand) {
			return apierror.BandNameTaken()
		}
		if appErr := asAppError(err); appErr != nil {
			return appErr
		}
		return apierror.InternalError("inscription")
	}

	RespondCreated(w, response)
	return nil
}

func (h UserHandler) Login(w http.ResponseWriter, r *http.Request) error {
	payload, err := DecodeJSON[service.LoginPayload](r)
	if err != nil {
		return err
	}

	response, err := h.UserService.Login(r.Context(), payload)
	if err != nil {
		return apierror.InvalidCredentials()
	}

	RespondOK(w, response)
	return nil
}

func (h UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) error {
	userID, err := GetUserID(r)
	if err != nil {
		return err
	}

	payload, err := DecodeJSON[service.UpdatePasswordPayload](r)
	if err != nil {
		return err
	}

	if payload.NewPassword == "" {
		return apierror.InvalidRequest("Le nouveau mot de passe ne peut pas être vide.")
	}

	if err := h.UserService.UpdatePassword(r.Context(), userID, payload); err != nil {
		if err.Error() == "invalid current password" {
			return apierror.WrongCurrentPassword()
		}
		if appErr := asAppError(err); appErr != nil {
			return appErr
		}
		return apierror.InternalError("mise à jour du mot de passe")
	}

	RespondOK(w, map[string]string{"message": "Mot de passe mis à jour avec succès."})
	return nil
}

func (h UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query().Get("q")

	users, err := h.UserService.SearchUsers(r.Context(), query)
	if err != nil {
		return apierror.InternalError("recherche d'utilisateurs")
	}
	if users == nil {
		users = make([]model.User, 0)
	}

	RespondOK(w, users)
	return nil
}
