package handler

import (
	"encoding/json"
	"net/http"
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
		writeError(w, "Could not identify band from token", http.StatusInternalServerError)
		return
	}

	var payload service.CreateSetlistPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	setlist, err := h.SetlistService.Create(r.Context(), payload, bandID)
	if err != nil {
		writeError(w, "Failed to create setlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(setlist)
}

func (h SetlistHandler) GetSetlists(w http.ResponseWriter, r *http.Request) {
	bandID, ok := r.Context().Value(middleware.BandIDKey).(int)
	if !ok {
		writeError(w, "Could not identify band from token", http.StatusInternalServerError)
		return
	}

	setlists, err := h.SetlistService.GetAllForBand(r.Context(), bandID)
	if err != nil {
		writeError(w, "Failed to retrieve setlists", http.StatusInternalServerError)
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
		writeError(w, "Could not identify band from token", http.StatusInternalServerError)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, "Invalid setlist ID", http.StatusBadRequest)
		return
	}

	details, err := h.SetlistService.GetDetails(r.Context(), id, bandID)
	if err != nil {
		writeError(w, "Setlist not found", http.StatusNotFound)
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
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	item, err := h.SetlistService.AddItem(r.Context(), setlistID, payload)
	if err != nil {
		writeError(w, "Failed to add item to setlist", http.StatusInternalServerError)
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
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.SetlistService.UpdateOrder(r.Context(), setlistID, payload); err != nil {
		writeError(w, "Failed to update setlist order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order updated successfully"})
}

func (h SetlistHandler) UpdateSetlist(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	id, _ := strconv.Atoi(r.PathValue("id"))

	var payload service.UpdateSetlistPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.SetlistService.Update(r.Context(), id, bandID, payload); err != nil {
		writeError(w, "Failed to update setlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h SetlistHandler) DeleteSetlist(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := h.SetlistService.Delete(r.Context(), id, bandID); err != nil {
		writeError(w, "Failed to delete setlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h SetlistHandler) DeleteSetlistItem(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	setlistID, _ := strconv.Atoi(r.PathValue("id"))
	itemID, _ := strconv.Atoi(r.PathValue("itemId"))

	if err := h.SetlistService.DeleteItem(r.Context(), itemID, setlistID, bandID); err != nil {
		writeError(w, "Failed to delete item from setlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
