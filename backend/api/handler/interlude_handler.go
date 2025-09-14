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

func (h InterludeHandler) GetInterlude(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	id, _ := strconv.Atoi(r.PathValue("id"))

	interlude, err := h.InterludeService.GetByID(r.Context(), id, bandID)
	if err != nil {
		writeError(w, "Interlude not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interlude)
}

func (h InterludeHandler) UpdateInterlude(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	id, _ := strconv.Atoi(r.PathValue("id"))

	var payload service.CreateInterludePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.InterludeService.Update(r.Context(), id, bandID, payload); err != nil {
		writeError(w, "Failed to update interlude", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h InterludeHandler) DeleteInterlude(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := h.InterludeService.Delete(r.Context(), id, bandID); err != nil {
		writeError(w, "Failed to delete interlude", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
