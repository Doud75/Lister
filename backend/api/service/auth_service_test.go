package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"setlist/api/repository/mocks"

	"go.uber.org/mock/gomock"
)

func TestAuthService_RevokeRefreshToken(t *testing.T) {
	ctx := context.Background()
	const refreshToken = "some-refresh-token"

	t.Run("revokes token owned by the user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockRefreshTokenRepository(ctrl)
		svc := AuthService{RefreshTokenRepo: mockRepo}

		mockRepo.EXPECT().
			FindRefreshToken(ctx, gomock.Any()).
			Return(1, time.Now().Add(time.Hour), nil)
		mockRepo.EXPECT().
			DeleteRefreshToken(ctx, gomock.Any()).
			Return(nil)

		if err := svc.RevokeRefreshToken(ctx, refreshToken, 1); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("refuses to revoke a token owned by another user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockRefreshTokenRepository(ctrl)
		svc := AuthService{RefreshTokenRepo: mockRepo}

		mockRepo.EXPECT().
			FindRefreshToken(ctx, gomock.Any()).
			Return(999, time.Now().Add(time.Hour), nil)
		// DeleteRefreshToken must never be called; the mock fails the test otherwise.

		err := svc.RevokeRefreshToken(ctx, refreshToken, 1)
		if err == nil {
			t.Fatal("expected an error when token belongs to another user")
		}
	})

	t.Run("propagates lookup error without deleting", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockRefreshTokenRepository(ctrl)
		svc := AuthService{RefreshTokenRepo: mockRepo}

		mockRepo.EXPECT().
			FindRefreshToken(ctx, gomock.Any()).
			Return(0, time.Time{}, errors.New("not found"))

		if err := svc.RevokeRefreshToken(ctx, refreshToken, 1); err == nil {
			t.Fatal("expected the lookup error to be propagated")
		}
	})
}
