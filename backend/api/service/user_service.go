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

type UserService struct {
	UserRepo         repository.UserRepository
	RefreshTokenRepo repository.RefreshTokenRepository
	JWTSecret        string
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
	Username string  `json:"username"`
	Password *string `json:"password"`
}

type AuthResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
	Bands        []model.Band `json:"bands"`
}

func (s UserService) Signup(ctx context.Context, payload AuthPayload) (*AuthResponse, error) {
	if err := ValidateUsername(payload.Username); err != nil {
		return nil, err
	}
	if err := ValidatePassword(payload.Password); err != nil {
		return nil, err
	}

	payload.Username = SanitizeString(payload.Username)
	payload.BandName = SanitizeString(payload.BandName)

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

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	tokenHash, err := auth.HashRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(auth.RefreshTokenDuration)
	err = s.RefreshTokenRepo.ReplaceUserRefreshToken(ctx, user.ID, tokenHash, expiresAt)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		Bands:        []model.Band{band},
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

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	tokenHash, err := auth.HashRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(auth.RefreshTokenDuration)
	err = s.RefreshTokenRepo.ReplaceUserRefreshToken(ctx, user.ID, tokenHash, expiresAt)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		Bands:        bands,
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

	if err := ValidatePassword(payload.NewPassword); err != nil {
		return err
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
	if payload.Username == "" {
		return model.User{}, errors.New("username is required")
	}
	if err := ValidateUsername(payload.Username); err != nil {
		return model.User{}, err
	}
	payload.Username = SanitizeString(payload.Username)

	existingUser, err := s.UserRepo.FindUserByUsername(ctx, payload.Username)
	if err == nil {
		isMember, err := s.UserRepo.IsUserInBand(ctx, existingUser.ID, bandID)
		if err != nil {
			return model.User{}, err
		}
		if isMember {
			return model.User{}, errors.New("user is already a member of this band")
		}

		err = s.UserRepo.AddUserToBand(ctx, existingUser.ID, bandID, "member")
		if err != nil {
			return model.User{}, err
		}
		return existingUser, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, err
	}

	if payload.Password == nil || *payload.Password == "" {
		return model.User{}, errors.New("user not found, and password is required to create a new one")
	}
	if err := ValidatePassword(*payload.Password); err != nil {
		return model.User{}, err
	}

	hashedPassword, err := auth.HashPassword(*payload.Password)
	if err != nil {
		return model.User{}, err
	}

	newUser, err := s.UserRepo.CreateUserAndAddToBand(ctx, bandID, payload.Username, hashedPassword, "member")
	if err != nil {
		return model.User{}, err
	}

	return newUser, nil
}

func (s UserService) RemoveMember(ctx context.Context, bandID int, userID int) error {
	return s.UserRepo.RemoveUserFromBand(ctx, bandID, userID)
}

func (s UserService) SearchUsers(ctx context.Context, query string) ([]model.User, error) {
	if len(query) < 3 {
		return []model.User{}, nil
	}
	return s.UserRepo.SearchUsersByUsername(ctx, query)
}
