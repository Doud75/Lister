package service

import (
	"context"
	"errors"
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

func (s UserService) Signup(ctx context.Context, payload AuthPayload) (string, error) {
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return "", err
	}

	user, err := s.UserRepo.CreateBandAndUser(ctx, payload.BandName, payload.Username, hashedPassword)
	if err != nil {
		return "", err
	}

	return auth.GenerateJWT(s.JWTSecret, user.ID, user.BandID)
}

func (s UserService) Join(ctx context.Context, payload AuthPayload) (string, error) {
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return "", err
	}

	user, err := s.UserRepo.CreateUserForExistingBand(ctx, payload.BandName, payload.Username, hashedPassword)
	if err != nil {
		return "", err
	}

	return auth.GenerateJWT(s.JWTSecret, user.ID, user.BandID)
}

func (s UserService) Login(ctx context.Context, payload LoginPayload) (string, error) {
	user, err := s.UserRepo.FindUserByUsername(ctx, payload.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !auth.CheckPasswordHash(payload.Password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	return auth.GenerateJWT(s.JWTSecret, user.ID, user.BandID)
}
