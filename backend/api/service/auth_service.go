package service

import (
	"context"
	"errors"
	"setlist/api/model"
	"setlist/api/repository"
	"setlist/auth"
	"time"

	"github.com/jackc/pgx/v5"
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
	tokenHash, err := auth.HashRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userID, expiresAt, err := s.RefreshTokenRepo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("refresh token not found")
		}
		return nil, err
	}

	if time.Now().After(expiresAt) {
		s.RefreshTokenRepo.DeleteRefreshToken(ctx, tokenHash)
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
	err = s.RefreshTokenRepo.StoreRefreshToken(ctx, userID, newRefreshTokenHash, newExpiresAt)
	if err != nil {
		return nil, err
	}
	s.RefreshTokenRepo.DeleteRefreshToken(ctx, tokenHash)
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
