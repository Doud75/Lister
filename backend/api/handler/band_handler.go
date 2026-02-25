package handler

import (
	"net/http"
	"setlist/api/apierror"
	"setlist/api/model"
	"setlist/api/service"
)

type BandHandler struct {
	UserService service.UserService
}

func (h BandHandler) GetMembers(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	members, err := h.UserService.GetBandMembers(r.Context(), bandID)
	if err != nil {
		return apierror.InternalError("récupération des membres")
	}
	if members == nil {
		members = make([]model.BandMember, 0)
	}

	RespondOK(w, members)
	return nil
}

func (h BandHandler) InviteMember(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	payload, err := DecodeJSON[service.InviteMemberPayload](r)
	if err != nil {
		return err
	}

	user, err := h.UserService.InviteMember(r.Context(), bandID, payload)
	if err != nil {
		if appErr := asAppError(err); appErr != nil {
			return appErr
		}
		return apierror.NewUserError(apierror.ErrInvalidRequest, err.Error(), http.StatusBadRequest)
	}

	RespondCreated(w, user)
	return nil
}

func (h BandHandler) RemoveMember(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	userID, err := GetIntParam(r, "userId")
	if err != nil {
		return apierror.InvalidRequest("Identifiant utilisateur invalide.")
	}

	if err := h.UserService.RemoveMember(r.Context(), bandID, userID); err != nil {
		return apierror.NewUserError(apierror.ErrInvalidRequest, err.Error(), http.StatusBadRequest)
	}

	RespondNoContent(w)
	return nil
}
