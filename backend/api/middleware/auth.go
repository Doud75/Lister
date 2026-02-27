package middleware

import (
	"context"
	"net/http"
	"setlist/api/repository"
	"setlist/auth"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const (
	UserIDKey ContextKey = "userID"
	BandIDKey ContextKey = "bandID"
)

func JWTAuth(jwtSecret string, userRepo repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claimsVal, ok := validateJWT(w, r, jwtSecret)
			if !ok {
				return
			}

			bandIDStr := r.Header.Get("X-Band-ID")
			if bandIDStr == "" {
				http.Error(w, "Missing X-Band-ID header", http.StatusBadRequest)
				return
			}
			bandID, err := strconv.Atoi(bandIDStr)
			if err != nil {
				http.Error(w, "Invalid X-Band-ID header", http.StatusBadRequest)
				return
			}

			isMember, err := userRepo.IsUserInBand(r.Context(), claimsVal.UserID, bandID)
			if err != nil {
				http.Error(w, "Error verifying band membership", http.StatusInternalServerError)
				return
			}
			if !isMember {
				http.Error(w, "User is not a member of this band", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claimsVal.UserID)
			ctx = context.WithValue(ctx, BandIDKey, bandID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func JWTAuthUserOnly(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claimsVal, ok := validateJWT(w, r, jwtSecret)
			if !ok {
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claimsVal.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateJWT(w http.ResponseWriter, r *http.Request, jwtSecret string) (*auth.JWTClaims, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing authorization header", http.StatusUnauthorized)
		return nil, false
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return nil, false
	}

	claims := &auth.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return nil, false
	}

	return claims, true
}