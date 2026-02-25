package handler

import (
	"net/http"
	"setlist/api/apierror"
	"setlist/api/repository"
)

type InfoHandler struct {
	InfoRepo repository.InfoRepository
	UserRepo repository.UserRepository
}

func (h InfoHandler) GetCurrentUserInfo(w http.ResponseWriter, r *http.Request) error {
	userID, err := GetUserID(r)
	if err != nil {
		return err
	}
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	user, err := h.InfoRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		return apierror.NotFound("Utilisateur")
	}

	band, err := h.InfoRepo.GetBandByID(r.Context(), bandID)
	if err != nil {
		return apierror.NotFound("Groupe")
	}

	role, err := h.UserRepo.GetUserRoleInBand(r.Context(), userID, bandID)
	if err != nil {
		return apierror.NotFound("Rôle utilisateur")
	}

	RespondOK(w, map[string]string{
		"username":  user.Username,
		"band_name": band.Name,
		"role":      role,
	})
	return nil
}
