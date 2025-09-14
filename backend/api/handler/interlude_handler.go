package handler

import (
	"encoding/json"
	"net/http"
	"setlist/api/middleware"
	"setlist/api/service"
	"strconv"
)

type InterludeHandler struct {
	InterludeService service.InterludeService
}

func (h InterludeHandler) CreateInterlude(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)

	var payload service.CreateInterludePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdInterlude, err := h.InterludeService.Create(r.Context(), payload, bandID)
	if err != nil {
		writeError(w, "Failed to create interlude", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdInterlude)
}

func (h InterludeHandler) GetInterludes(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	interludes, err := h.InterludeService.GetAllForBand(r.Context(), bandID)
	if err != nil {
		writeError(w, "Failed to retrieve interludes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interludes)
}

func (h InterludeHandler) UpdateInterlude(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	var payload service.UpdateInterludePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedInterlude, err := h.InterludeService.Update(r.Context(), id, bandID, payload)
	if err != nil {
		writeError(w, "Failed to update interlude", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedInterlude)
}
