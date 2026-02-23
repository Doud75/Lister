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
	"strconv"
)

type BandHandler struct {
	UserService service.UserService
}

func (h BandHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)

	members, err := h.UserService.GetBandMembers(r.Context(), bandID)
	if err != nil {
		writeAppError(w, apierror.InternalError("récupération des membres"))
		return
	}
	if members == nil {
		members = make([]model.BandMember, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func (h BandHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)

	var payload service.InviteMemberPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeAppError(w, apierror.InvalidRequest("Corps de la requête invalide."))
		return
	}

	user, err := h.UserService.InviteMember(r.Context(), bandID, payload)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateUsername) {
			writeAppError(w, apierror.UsernameTaken())
			return
		}
		var appErr *apierror.AppError
		if errors.As(err, &appErr) {
			writeAppError(w, appErr)
			return
		}
		writeAppError(w, apierror.NewUserError(apierror.ErrInvalidRequest, err.Error(), http.StatusBadRequest))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h BandHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	userIDStr := r.PathValue("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeAppError(w, apierror.InvalidRequest("Identifiant utilisateur invalide."))
		return
	}

	err = h.UserService.RemoveMember(r.Context(), bandID, userID)
	if err != nil {
		writeAppError(w, apierror.NewUserError(apierror.ErrInvalidRequest, err.Error(), http.StatusBadRequest))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
