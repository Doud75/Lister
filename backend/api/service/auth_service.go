package service

import (
	"context"
	"errors"
	"setlist/api/model"
	"setlist/api/repository"
	"setlist/auth"
	"time"
)

type AuthService struct {
	UserRepo         repository.UserRepository
	RefreshTokenRepo repository.RefreshTokenRepository
	JWTSecret        string
}

type RefreshTokenResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
	Bands        []model.Band `json:"bands"`
}

func (s AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	var userID int
	var expiresAt time.Time
	var matchedHash string

	query := `SELECT user_id, token_hash, expires_at FROM refresh_tokens WHERE expires_at > NOW()`
	rows, err := s.RefreshTokenRepo.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var uid int
		var hash string
		var exp time.Time
		if err := rows.Scan(&uid, &hash, &exp); err != nil {
			continue
		}

		if auth.VerifyRefreshToken(refreshToken, hash) {
			userID = uid
			expiresAt = exp
			matchedHash = hash
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("refresh token not found")
	}

	if time.Now().After(expiresAt) {
		s.RefreshTokenRepo.DeleteRefreshToken(ctx, matchedHash)
		return nil, errors.New("refresh token expired")
	}

	newAccessToken, err := auth.GenerateJWT(s.JWTSecret, userID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newRefreshTokenHash, err := auth.HashRefreshToken(newRefreshToken)
	if err != nil {
		return nil, err
	}

	newExpiresAt := time.Now().Add(auth.RefreshTokenDuration)
	err = s.RefreshTokenRepo.ReplaceUserRefreshToken(ctx, userID, newRefreshTokenHash, newExpiresAt)
	if err != nil {
		return nil, err
	}

	bands, err := s.UserRepo.FindBandsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
		Bands:        bands,
	}, nil
}

func (s AuthService) RevokeRefreshToken(ctx context.Context, refreshToken string, userID int) error {
	hashes, err := s.RefreshTokenRepo.GetUserTokenHashes(ctx, userID)
	if err != nil {
		return err
	}

	for _, hash := range hashes {
		if auth.VerifyRefreshToken(refreshToken, hash) {
			return s.RefreshTokenRepo.DeleteRefreshToken(ctx, hash)
		}
	}

	return nil
}

func (s AuthService) RevokeAllUserTokens(ctx context.Context, userID int) error {
	return s.RefreshTokenRepo.DeleteAllUserTokens(ctx, userID)
}
