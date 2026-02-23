package handler

import (
	"encoding/json"
	"net/http"
	"setlist/api/apierror"
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
		writeAppError(w, apierror.NewServerError(apierror.ErrInternal, "Impossible d'identifier l'utilisateur depuis le token."))
		return
	}

	user, err := h.InfoRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		writeAppError(w, apierror.NotFound("Utilisateur"))
		return
	}

	band, err := h.InfoRepo.GetBandByID(r.Context(), bandID)
	if err != nil {
		writeAppError(w, apierror.NotFound("Groupe"))
		return
	}

	role, err := h.UserRepo.GetUserRoleInBand(r.Context(), userID, bandID)
	if err != nil {
		writeAppError(w, apierror.NotFound("Rôle utilisateur"))
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
