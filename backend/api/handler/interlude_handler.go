package handler

import (
	"net/http"
	"setlist/api/apierror"
	"setlist/api/service"
)

type InterludeHandler struct {
	InterludeService service.InterludeService
}

func (h InterludeHandler) CreateInterlude(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	payload, err := DecodeJSON[service.CreateInterludePayload](r)
	if err != nil {
		return err
	}

	createdInterlude, err := h.InterludeService.Create(r.Context(), payload, bandID)
	if err != nil {
		return apierror.InternalError("création d'interlude")
	}

	RespondCreated(w, createdInterlude)
	return nil
}

func (h InterludeHandler) GetInterludes(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	interludes, err := h.InterludeService.GetAllForBand(r.Context(), bandID)
	if err != nil {
		return apierror.InternalError("récupération des interludes")
	}

	RespondOK(w, interludes)
	return nil
}

func (h InterludeHandler) UpdateInterlude(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	id, err := GetIntParam(r, "id")
	if err != nil {
		return apierror.InvalidRequest("Identifiant d'interlude invalide.")
	}

	payload, err := DecodeJSON[service.UpdateInterludePayload](r)
	if err != nil {
		return err
	}

	updatedInterlude, err := h.InterludeService.Update(r.Context(), id, bandID, payload)
	if err != nil {
		return apierror.NotFound("Interlude")
	}

	RespondOK(w, updatedInterlude)
	return nil
}
