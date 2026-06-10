package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"setlist/api/apierror"
	"setlist/api/repository"
	"setlist/cache"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

const profileCacheTTL = 15 * time.Minute

type InfoHandler struct {
	InfoRepo repository.InfoRepository
	UserRepo repository.UserRepository
	Cache    *redis.Client
}

// notFoundOrInternal returns a 404 for the given entity only when the error is
// a driver-level "no rows"; anything else is reported as an internal error.
func notFoundOrInternal(err error, entity, operation string) error {
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		return apierror.NotFound(entity)
	}
	return apierror.InternalError(operation)
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
		return notFoundOrInternal(err, "Utilisateur", "récupération de l'utilisateur")
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
		return notFoundOrInternal(err, "Groupe", "récupération du groupe")
	}

	role, err := h.UserRepo.GetUserRoleInBand(r.Context(), userID, bandID)
	if err != nil {
		return notFoundOrInternal(err, "Rôle utilisateur", "récupération du rôle")
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
