package handler

import (
	"encoding/json"
	"net/http"
	"setlist/api/apierror"
	"setlist/api/middleware"
	"setlist/api/model"
	"setlist/api/service"
	"strconv"
)

type SetlistHandler struct {
	SetlistService service.SetlistService
}

func (h SetlistHandler) CreateSetlist(w http.ResponseWriter, r *http.Request) {
	bandID, ok := r.Context().Value(middleware.BandIDKey).(int)
	if !ok {
		writeAppError(w, apierror.NewServerError(apierror.ErrInternal, "Impossible d'identifier le groupe depuis le token."))
		return
	}

	var payload service.CreateSetlistPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeAppError(w, apierror.InvalidRequest("Corps de la requête invalide."))
		return
	}

	setlist, err := h.SetlistService.Create(r.Context(), payload, bandID)
	if err != nil {
		writeAppError(w, apierror.InternalError("création de setlist"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(setlist)
}

func (h SetlistHandler) UpdateSetlist(w http.ResponseWriter, r *http.Request) {
	bandID, ok := r.Context().Value(middleware.BandIDKey).(int)
	if !ok {
		writeAppError(w, apierror.NewServerError(apierror.ErrInternal, "Impossible d'identifier le groupe depuis le token."))
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeAppError(w, apierror.InvalidRequest("Identifiant de setlist invalide."))
		return
	}

	var payload service.UpdateSetlistPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeAppError(w, apierror.InvalidRequest("Corps de la requête invalide."))
		return
	}

	setlist, err := h.SetlistService.Update(r.Context(), id, bandID, payload)
	if err != nil {
		writeAppError(w, apierror.InternalError("mise à jour de setlist"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(setlist)
}

func (h SetlistHandler) DeleteSetlist(w http.ResponseWriter, r *http.Request) {
	bandID, ok := r.Context().Value(middleware.BandIDKey).(int)
	if !ok {
		writeAppError(w, apierror.NewServerError(apierror.ErrInternal, "Impossible d'identifier le groupe depuis le token."))
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeAppError(w, apierror.InvalidRequest("Identifiant de setlist invalide."))
		return
	}

	err = h.SetlistService.Delete(r.Context(), id, bandID)
	if err != nil {
		writeAppError(w, apierror.InternalError("suppression de setlist"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h SetlistHandler) GetSetlists(w http.ResponseWriter, r *http.Request) {
	bandID, ok := r.Context().Value(middleware.BandIDKey).(int)
	if !ok {
		writeAppError(w, apierror.NewServerError(apierror.ErrInternal, "Impossible d'identifier le groupe depuis le token."))
		return
	}

	setlists, err := h.SetlistService.GetAllForBand(r.Context(), bandID)
	if err != nil {
		writeAppError(w, apierror.InternalError("récupération des setlists"))
		return
	}

	if setlists == nil {
		setlists = make([]model.Setlist, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(setlists)
}

func (h SetlistHandler) GetSetlistDetails(w http.ResponseWriter, r *http.Request) {
	bandID, ok := r.Context().Value(middleware.BandIDKey).(int)
	if !ok {
		writeAppError(w, apierror.NewServerError(apierror.ErrInternal, "Impossible d'identifier le groupe depuis le token."))
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeAppError(w, apierror.InvalidRequest("Identifiant de setlist invalide."))
		return
	}

	details, err := h.SetlistService.GetDetails(r.Context(), id, bandID)
	if err != nil {
		writeAppError(w, apierror.NotFound("Setlist"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

func (h SetlistHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	setlistIDStr := r.PathValue("id")
	setlistID, _ := strconv.Atoi(setlistIDStr)

	var payload service.AddItemPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeAppError(w, apierror.InvalidRequest("Corps de la requête invalide."))
		return
	}

	item, err := h.SetlistService.AddItem(r.Context(), setlistID, payload)
	if err != nil {
		writeAppError(w, apierror.InternalError("ajout d'élément à la setlist"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h SetlistHandler) UpdateItemOrder(w http.ResponseWriter, r *http.Request) {
	setlistIDStr := r.PathValue("id")
	setlistID, _ := strconv.Atoi(setlistIDStr)

	var payload service.UpdateOrderPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeAppError(w, apierror.InvalidRequest("Corps de la requête invalide."))
		return
	}

	if err := h.SetlistService.UpdateOrder(r.Context(), setlistID, payload); err != nil {
		writeAppError(w, apierror.InternalError("mise à jour de l'ordre"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h SetlistHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	itemIDStr := r.PathValue("itemId")
	itemID, _ := strconv.Atoi(itemIDStr)

	var payload service.UpdateItemPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeAppError(w, apierror.InvalidRequest("Corps de la requête invalide."))
		return
	}

	item, err := h.SetlistService.UpdateItem(r.Context(), itemID, bandID, payload)
	if err != nil {
		writeAppError(w, apierror.InternalError("mise à jour d'élément"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h SetlistHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	itemIDStr := r.PathValue("itemId")
	itemID, _ := strconv.Atoi(itemIDStr)

	if err := h.SetlistService.DeleteItem(r.Context(), itemID, bandID); err != nil {
		writeAppError(w, apierror.InternalError("suppression d'élément"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h SetlistHandler) DuplicateSetlist(w http.ResponseWriter, r *http.Request) {
	bandID, ok := r.Context().Value(middleware.BandIDKey).(int)
	if !ok {
		writeAppError(w, apierror.NewServerError(apierror.ErrInternal, "Impossible d'identifier le groupe depuis le token."))
		return
	}

	idStr := r.PathValue("id")
	originalSetlistID, err := strconv.Atoi(idStr)
	if err != nil {
		writeAppError(w, apierror.InvalidRequest("Identifiant de setlist invalide."))
		return
	}

	var payload service.DuplicateSetlistPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeAppError(w, apierror.InvalidRequest("Corps de la requête invalide."))
		return
	}

	newSetlist, err := h.SetlistService.Duplicate(r.Context(), originalSetlistID, bandID, payload.Name, payload.Color)
	if err != nil {
		writeAppError(w, apierror.InternalError("duplication de setlist"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSetlist)
}
