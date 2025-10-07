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

type UpdatePasswordPayload struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type InviteMemberPayload struct {
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

func (s UserService) UpdatePassword(ctx context.Context, userID int, payload UpdatePasswordPayload) error {
	user, err := s.UserRepo.FindUserByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !auth.CheckPasswordHash(payload.CurrentPassword, user.PasswordHash) {
		return errors.New("invalid current password")
	}

	newHashedPassword, err := auth.HashPassword(payload.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	return s.UserRepo.UpdatePassword(ctx, userID, newHashedPassword)
}

func (s UserService) GetBandMembers(ctx context.Context, bandID int) ([]model.BandMember, error) {
	return s.UserRepo.GetMembersByBandID(ctx, bandID)
}

func (s UserService) InviteMember(ctx context.Context, bandID int, payload InviteMemberPayload) (model.User, error) {
	if payload.Username == "" || payload.Password == "" {
		return model.User{}, errors.New("username and password are required")
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return model.User{}, err
	}

	user, err := s.UserRepo.CreateUserAndAddToBand(ctx, bandID, payload.Username, hashedPassword, "member")
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s UserService) RemoveMember(ctx context.Context, bandID int, userID int) error {
	return s.UserRepo.RemoveUserFromBand(ctx, bandID, userID)
}
