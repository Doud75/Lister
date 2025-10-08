package handler

import (
	"encoding/json"
	"net/http"
	"setlist/api/middleware"
	"setlist/api/repository"
)

type InfoHandler struct {
	InfoRepo repository.InfoRepository
	UserRepo repository.UserRepository
}

func (h InfoHandler) GetCurrentUserInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	bandID, ok2 := r.Context().Value(middleware.BandIDKey).(int)
	if !ok || !ok2 {
		writeError(w, "Could not identify user from token", http.StatusInternalServerError)
		return
	}

	user, err := h.InfoRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		writeError(w, "User not found", http.StatusNotFound)
		return
	}

	band, err := h.InfoRepo.GetBandByID(r.Context(), bandID)
	if err != nil {
		writeError(w, "Band not found", http.StatusNotFound)
		return
	}

	role, err := h.UserRepo.GetUserRoleInBand(r.Context(), userID, bandID)
	if err != nil {
		writeError(w, "Could not determine user role", http.StatusNotFound)
		return
	}

	response := map[string]string{
		"username":  user.Username,
		"band_name": band.Name,
		"role":      role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
