package service

import (
	"context"
	"testing"
	"time"

	"setlist/api/model"
	"setlist/api/repository"
	"setlist/api/repository/mocks"

	"go.uber.org/mock/gomock"
)

func TestInvitationService_CreateInvitation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvRepo := mocks.NewMockInvitationRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	svc := InvitationService{
		InvitationRepo: mockInvRepo,
		UserRepo:       mockUserRepo,
	}

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockInvRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(nil)

		token, expiresAt, err := svc.CreateInvitation(ctx, 1, "member")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(token) != 32 {
			t.Errorf("expected 32-char hex token, got length %d", len(token))
		}
		if expiresAt.Before(time.Now()) {
			t.Errorf("expected future expiration time")
		}
	})
}

func TestInvitationService_GetInvitationDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvRepo := mocks.NewMockInvitationRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	svc := InvitationService{
		InvitationRepo: mockInvRepo,
		UserRepo:       mockUserRepo,
	}

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		detail := &model.InvitationDetail{
			Token:     "testtoken",
			Role:      "member",
			BandName:  "The Beatles",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		mockInvRepo.EXPECT().
			GetByToken(ctx, "testtoken").
			Return(detail, nil)

		res, err := svc.GetInvitationDetails(ctx, "testtoken")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.BandName != "The Beatles" {
			t.Errorf("expected The Beatles, got %s", res.BandName)
		}
	})

	t.Run("Expired", func(t *testing.T) {
		detail := &model.InvitationDetail{
			Token:     "testtoken",
			Role:      "member",
			BandName:  "The Beatles",
			ExpiresAt: time.Now().Add(-24 * time.Hour),
		}

		mockInvRepo.EXPECT().
			GetByToken(ctx, "testtoken").
			Return(detail, nil)

		_, err := svc.GetInvitationDetails(ctx, "testtoken")
		if err != repository.ErrInvitationExpired {
			t.Errorf("expected ErrInvitationExpired, got %v", err)
		}
	})
}

func TestInvitationService_AcceptInvitation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvRepo := mocks.NewMockInvitationRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	svc := InvitationService{
		InvitationRepo: mockInvRepo,
		UserRepo:       mockUserRepo,
	}

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		detail := &model.InvitationDetail{
			Token:     "testtoken",
			BandID:    10,
			Role:      "member",
			BandName:  "The Beatles",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		mockInvRepo.EXPECT().
			GetByToken(ctx, "testtoken").
			Return(detail, nil)

		mockUserRepo.EXPECT().
			IsUserInBand(ctx, 42, 10).
			Return(false, nil)

		mockUserRepo.EXPECT().
			AddUserToBand(ctx, 42, 10, "member").
			Return(nil)

		res, err := svc.AcceptInvitation(ctx, "testtoken", 42)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.BandID != 10 {
			t.Errorf("expected band ID 10, got %d", res.BandID)
		}
	})

	t.Run("AlreadyMember", func(t *testing.T) {
		detail := &model.InvitationDetail{
			Token:     "testtoken",
			BandID:    10,
			Role:      "member",
			BandName:  "The Beatles",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		mockInvRepo.EXPECT().
			GetByToken(ctx, "testtoken").
			Return(detail, nil)

		mockUserRepo.EXPECT().
			IsUserInBand(ctx, 42, 10).
			Return(true, nil)

		_, err := svc.AcceptInvitation(ctx, "testtoken", 42)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "conflict: user is already a member" {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}
