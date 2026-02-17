package service

import (
	"context"
	"errors"
	"testing"

	"setlist/api/model"
	"setlist/api/repository/mocks"

	"go.uber.org/mock/gomock"
)

func TestUserService_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRefreshTokenRepo := mocks.NewMockRefreshTokenRepository(ctrl)

	jwtSecret := "testsecret"
	svc := UserService{
		UserRepo:         mockUserRepo,
		RefreshTokenRepo: mockRefreshTokenRepo,
		JWTSecret:        jwtSecret,
	}

	ctx := context.Background()
	payload := AuthPayload{
		BandName: "Test Band",
		Username: "testuser",
		Password: "Password123!",
	}

	t.Run("Success", func(t *testing.T) {
		expectedUser := model.User{ID: 1, Username: "testuser"}
		expectedBand := model.Band{ID: 1, Name: "Test Band"}

		// Expect CreateBandAndUser
		// Password argument is hashed inside the service, so we use gomock.Any() to match the hashed password string
		mockUserRepo.EXPECT().
			CreateBandAndUser(ctx, payload.BandName, payload.Username, gomock.Any()).
			Return(expectedUser, expectedBand, nil)

		// Expect ReplaceUserRefreshToken
		// Token hash and expiry are generated inside service, so allow Any()
		mockRefreshTokenRepo.EXPECT().
			ReplaceUserRefreshToken(ctx, expectedUser.ID, gomock.Any(), gomock.Any()).
			Return(nil)

		resp, err := svc.Signup(ctx, payload)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("response should not be nil")
		}
		if resp.Token == "" {
			t.Error("token should not be empty")
		}
		if resp.RefreshToken == "" {
			t.Error("refresh token should not be empty")
		}
		if len(resp.Bands) != 1 {
			t.Errorf("expected 1 band, got %d", len(resp.Bands))
		}
		if resp.Bands[0].Name != expectedBand.Name {
			t.Errorf("expected band name %s, got %s", expectedBand.Name, resp.Bands[0].Name)
		}
	})

	t.Run("RepoError", func(t *testing.T) {
		repoErr := errors.New("database error")

		mockUserRepo.EXPECT().
			CreateBandAndUser(ctx, payload.BandName, payload.Username, gomock.Any()).
			Return(model.User{}, model.Band{}, repoErr)

		resp, err := svc.Signup(ctx, payload)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if resp != nil {
			t.Fatal("response should be nil on error")
		}
		if !errors.Is(err, repoErr) {
			t.Errorf("expected error %v, got %v", repoErr, err)
		}
	})
}
