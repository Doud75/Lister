package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"setlist/api/model"
	"setlist/api/repository"
	"time"
)

type InvitationService struct {
	InvitationRepo repository.InvitationRepository
	UserRepo       repository.UserRepository
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s InvitationService) CreateInvitation(ctx context.Context, bandID int, role string) (string, time.Time, error) {
	token, err := generateToken()
	if err != nil {
		return "", time.Time{}, err
	}

	if role == "" {
		role = "member"
	}

	expiresAt := time.Now().Add(72 * time.Hour)

	invitation := &model.Invitation{
		Token:     token,
		BandID:    bandID,
		Role:      role,
		ExpiresAt: expiresAt,
		MaxUses:   nil,
	}

	if err := s.InvitationRepo.Create(ctx, invitation); err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (s InvitationService) GetInvitationDetails(ctx context.Context, token string) (*model.InvitationDetail, error) {
	detail, err := s.InvitationRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if time.Now().After(detail.ExpiresAt) {
		return nil, repository.ErrInvitationExpired
	}

	return detail, nil
}

func (s InvitationService) AcceptInvitation(ctx context.Context, token string, userID int) (*model.InvitationDetail, error) {
	invitation, err := s.InvitationRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if time.Now().After(invitation.ExpiresAt) {
		return nil, repository.ErrInvitationExpired
	}

	isMember, err := s.UserRepo.IsUserInBand(ctx, userID, invitation.BandID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return invitation, errors.New("conflict: user is already a member")
	}

	err = s.UserRepo.AddUserToBand(ctx, userID, invitation.BandID, invitation.Role)
	if err != nil {
		return nil, err
	}

	return invitation, nil
}
