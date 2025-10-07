package middleware

import (
	"net/http"
	"setlist/api/repository"
)

func AdminOnly(userRepo repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(UserIDKey).(int)
			if !ok {
				http.Error(w, "User not identified", http.StatusInternalServerError)
				return
			}
			bandID, ok := r.Context().Value(BandIDKey).(int)
			if !ok {
				http.Error(w, "Band not identified", http.StatusInternalServerError)
				return
			}

			role, err := userRepo.GetUserRoleInBand(r.Context(), userID, bandID)
			if err != nil {
				http.Error(w, "Could not verify user role", http.StatusInternalServerError)
				return
			}

			if role != "admin" {
				http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
