package handler

import (
	"encoding/json"
	"net/http"
	"setlist/api/middleware"
	"setlist/api/model"
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
		writeError(w, "Failed to retrieve members", http.StatusInternalServerError)
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
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.InviteMember(r.Context(), bandID, payload)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
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
		writeError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.UserService.RemoveMember(r.Context(), bandID, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
