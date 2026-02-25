package handler

import (
	"log"
	"net/http"
	"setlist/api/apierror"
	"setlist/api/service"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	AuthService service.AuthService
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	payload, err := DecodeJSON[RefreshTokenRequest](r)
	if err != nil {
		return err
	}

	if payload.RefreshToken == "" {
		return apierror.InvalidRequest("Le token de rafraîchissement est requis.")
	}

	response, err := h.AuthService.RefreshAccessToken(r.Context(), payload.RefreshToken)
	if err != nil {
		log.Printf("[AUTH] Refresh token failed: %v", err)
		return apierror.InvalidRefreshToken()
	}

	RespondOK(w, response)
	return nil
}

func (h AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	payload, err := DecodeJSON[RefreshTokenRequest](r)
	if err != nil {
		return err
	}

	if payload.RefreshToken != "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondOK(w, map[string]string{"message": "Déconnexion réussie."})
			return nil
		}

		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		if tokenString != "" {
			token, _ := jwt.Parse(tokenString, nil)
			if token != nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if userID, ok := claims["user_id"].(float64); ok {
						if err := h.AuthService.RevokeRefreshToken(r.Context(), payload.RefreshToken, int(userID)); err != nil {
							log.Printf("[AUTH] Failed to revoke refresh token: %v", err)
						} else {
							log.Printf("[AUTH] Refresh token revoked successfully for user %d", int(userID))
						}
					}
				}
			}
		}
	}

	RespondOK(w, map[string]string{"message": "Déconnexion réussie."})
	return nil
}
