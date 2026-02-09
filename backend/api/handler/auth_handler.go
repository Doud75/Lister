package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"setlist/api/service"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	AuthService service.AuthService
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var payload RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.RefreshToken == "" {
		writeError(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	response, err := h.AuthService.RefreshAccessToken(r.Context(), payload.RefreshToken)
	if err != nil {
		log.Printf("[AUTH] Refresh token failed: %v", err)
		writeError(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	log.Printf("[AUTH] Token refreshed successfully for user")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var payload RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.RefreshToken != "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
			return
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
