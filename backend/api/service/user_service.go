package service

import (
	"context"
	"errors"
	"setlist/api/model"
	"setlist/api/repository"
	"setlist/auth"
)

type UserService struct {
	UserRepo  repository.UserRepository
	JWTSecret string
}

type AuthPayload struct {
	BandName string `json:"band_name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	Bands []model.Band `json:"bands"`
}

func (s UserService) Signup(ctx context.Context, payload AuthPayload) (*AuthResponse, error) {
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	user, band, err := s.UserRepo.CreateBandAndUser(ctx, payload.BandName, payload.Username, hashedPassword)
	if err != nil {
		return nil, err
	}

	token, err := auth.GenerateJWT(s.JWTSecret, user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		Bands: []model.Band{band},
	}, nil
}

func (s UserService) Join(ctx context.Context, payload AuthPayload) (*AuthResponse, error) {
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	user, band, err := s.UserRepo.CreateUserForExistingBand(ctx, payload.BandName, payload.Username, hashedPassword)
	if err != nil {
		return nil, err
	}

	token, err := auth.GenerateJWT(s.JWTSecret, user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		Bands: []model.Band{band},
	}, nil
}

func (s UserService) Login(ctx context.Context, payload LoginPayload) (*AuthResponse, error) {
	user, err := s.UserRepo.FindUserByUsername(ctx, payload.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !auth.CheckPasswordHash(payload.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	bands, err := s.UserRepo.FindBandsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if len(bands) == 0 {
		return nil, errors.New("user is not part of any band")
	}

	token, err := auth.GenerateJWT(s.JWTSecret, user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		Bands: bands,
	}, nil
}