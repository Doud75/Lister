package handler

import (
	"encoding/json"
	"net/http"
	"setlist/api/apierror"
	"setlist/api/repository"
	"setlist/cache"
	"time"

	"github.com/redis/go-redis/v9"
)

const profileCacheTTL = 15 * time.Minute

type InfoHandler struct {
	InfoRepo repository.InfoRepository
	UserRepo repository.UserRepository
	Cache    *redis.Client
}

type userInfoResponse struct {
	Username string `json:"username"`
	BandName string `json:"band_name"`
	Role     string `json:"role"`
}

func (h InfoHandler) GetCurrentUserInfo(w http.ResponseWriter, r *http.Request) error {
	userID, err := GetUserID(r)
	if err != nil {
		return err
	}

	bandID, hasBand := GetOptionalBandID(r)

	user, err := h.InfoRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		return apierror.NotFound("Utilisateur")
	}

	if !hasBand {
		RespondOK(w, userInfoResponse{
			Username: user.Username,
			BandName: "",
			Role:     "",
		})
		return nil
	}

	cacheKey := cache.ProfileKey(userID, bandID)
	if data, ok := cache.Get(r.Context(), h.Cache, cacheKey); ok {
		var info userInfoResponse
		if err := json.Unmarshal([]byte(data), &info); err == nil {
			RespondOK(w, info)
			return nil
		}
	}

	band, err := h.InfoRepo.GetBandByID(r.Context(), bandID)
	if err != nil {
		return apierror.NotFound("Groupe")
	}

	role, err := h.UserRepo.GetUserRoleInBand(r.Context(), userID, bandID)
	if err != nil {
		return apierror.NotFound("Rôle utilisateur")
	}

	info := userInfoResponse{
		Username: user.Username,
		BandName: band.Name,
		Role:     role,
	}

	if data, err := json.Marshal(info); err == nil {
		cache.Set(r.Context(), h.Cache, cacheKey, string(data), profileCacheTTL)
	}

	RespondOK(w, info)
	return nil
}
