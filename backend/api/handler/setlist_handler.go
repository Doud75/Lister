package handler

import (
	"net/http"
	"setlist/api/apierror"
	"setlist/api/model"
	"setlist/api/service"
)

type SetlistHandler struct {
	SetlistService service.SetlistService
}

func (h SetlistHandler) CreateSetlist(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	payload, err := DecodeJSON[service.CreateSetlistPayload](r)
	if err != nil {
		return err
	}

	setlist, err := h.SetlistService.Create(r.Context(), payload, bandID)
	if err != nil {
		return apierror.InternalError("création de setlist")
	}

	RespondCreated(w, setlist)
	return nil
}

func (h SetlistHandler) UpdateSetlist(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	id, err := GetIntParam(r, "id")
	if err != nil {
		return apierror.InvalidRequest("Identifiant de setlist invalide.")
	}

	payload, err := DecodeJSON[service.UpdateSetlistPayload](r)
	if err != nil {
		return err
	}

	setlist, err := h.SetlistService.Update(r.Context(), id, bandID, payload)
	if err != nil {
		return apierror.InternalError("mise à jour de setlist")
	}

	RespondOK(w, setlist)
	return nil
}

func (h SetlistHandler) DeleteSetlist(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	id, err := GetIntParam(r, "id")
	if err != nil {
		return apierror.InvalidRequest("Identifiant de setlist invalide.")
	}

	if err := h.SetlistService.Delete(r.Context(), id, bandID); err != nil {
		return apierror.InternalError("suppression de setlist")
	}

	RespondNoContent(w)
	return nil
}

func (h SetlistHandler) GetSetlists(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	setlists, err := h.SetlistService.GetAllForBand(r.Context(), bandID)
	if err != nil {
		return apierror.InternalError("récupération des setlists")
	}
	if setlists == nil {
		setlists = make([]model.Setlist, 0)
	}

	RespondOK(w, setlists)
	return nil
}

func (h SetlistHandler) GetSetlistDetails(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	id, err := GetIntParam(r, "id")
	if err != nil {
		return apierror.InvalidRequest("Identifiant de setlist invalide.")
	}

	details, err := h.SetlistService.GetDetails(r.Context(), id, bandID)
	if err != nil {
		return apierror.NotFound("Setlist")
	}

	RespondOK(w, details)
	return nil
}

func (h SetlistHandler) AddItem(w http.ResponseWriter, r *http.Request) error {
	setlistID, err := GetIntParam(r, "id")
	if err != nil {
		return apierror.InvalidRequest("Identifiant de setlist invalide.")
	}

	payload, err := DecodeJSON[service.AddItemPayload](r)
	if err != nil {
		return err
	}

	item, err := h.SetlistService.AddItem(r.Context(), setlistID, payload)
	if err != nil {
		return apierror.InternalError("ajout d'élément à la setlist")
	}

	RespondCreated(w, item)
	return nil
}

func (h SetlistHandler) UpdateItemOrder(w http.ResponseWriter, r *http.Request) error {
	setlistID, err := GetIntParam(r, "id")
	if err != nil {
		return apierror.InvalidRequest("Identifiant de setlist invalide.")
	}

	payload, err := DecodeJSON[service.UpdateOrderPayload](r)
	if err != nil {
		return err
	}

	if err := h.SetlistService.UpdateOrder(r.Context(), setlistID, payload); err != nil {
		return apierror.InternalError("mise à jour de l'ordre")
	}

	RespondNoContent(w)
	return nil
}

func (h SetlistHandler) UpdateItem(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	itemID, err := GetIntParam(r, "itemId")
	if err != nil {
		return apierror.InvalidRequest("Identifiant d'élément invalide.")
	}

	payload, err := DecodeJSON[service.UpdateItemPayload](r)
	if err != nil {
		return err
	}

	item, err := h.SetlistService.UpdateItem(r.Context(), itemID, bandID, payload)
	if err != nil {
		return apierror.InternalError("mise à jour d'élément")
	}

	RespondOK(w, item)
	return nil
}

func (h SetlistHandler) DeleteItem(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	itemID, err := GetIntParam(r, "itemId")
	if err != nil {
		return apierror.InvalidRequest("Identifiant d'élément invalide.")
	}

	if err := h.SetlistService.DeleteItem(r.Context(), itemID, bandID); err != nil {
		return apierror.InternalError("suppression d'élément")
	}

	RespondNoContent(w)
	return nil
}

func (h SetlistHandler) DuplicateSetlist(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	originalSetlistID, err := GetIntParam(r, "id")
	if err != nil {
		return apierror.InvalidRequest("Identifiant de setlist invalide.")
	}

	payload, err := DecodeJSON[service.DuplicateSetlistPayload](r)
	if err != nil {
		return err
	}

	newSetlist, err := h.SetlistService.Duplicate(r.Context(), originalSetlistID, bandID, payload.Name, payload.Color)
	if err != nil {
		return apierror.InternalError("duplication de setlist")
	}

	RespondCreated(w, newSetlist)
	return nil
}
