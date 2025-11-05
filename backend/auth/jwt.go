package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	BandName string `json:"band_name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type UserForToken struct {
	ID       int
	Username string
	BandName string
	Role     string
}

func GenerateJWT(secretKey string, user UserForToken) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		BandName: user.BandName,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}