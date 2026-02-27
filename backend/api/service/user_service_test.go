package service

import (
	"context"
	"errors"
	"testing"

	"setlist/api/model"
	"setlist/api/repository/mocks"
	"setlist/auth"

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

		mockUserRepo.EXPECT().
			CreateBandAndUser(ctx, payload.BandName, payload.Username, gomock.Any()).
			Return(expectedUser, expectedBand, nil)
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

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRefreshTokenRepo := mocks.NewMockRefreshTokenRepository(ctrl)

	svc := UserService{
		UserRepo:         mockUserRepo,
		RefreshTokenRepo: mockRefreshTokenRepo,
		JWTSecret:        "testsecret",
	}

	ctx := context.Background()
	payload := LoginPayload{Username: "testuser", Password: "Password123!"}

	t.Run("LoginWithBands", func(t *testing.T) {
		hashedPw, _ := hashPasswordForTest("Password123!")
		expectedUser := model.User{ID: 1, Username: "testuser", PasswordHash: hashedPw}
		expectedBands := []model.Band{{ID: 1, Name: "Test Band"}}

		mockUserRepo.EXPECT().
			FindUserByUsername(ctx, payload.Username).
			Return(expectedUser, nil)
		mockUserRepo.EXPECT().
			FindBandsByUserID(ctx, expectedUser.ID).
			Return(expectedBands, nil)
		mockRefreshTokenRepo.EXPECT().
			ReplaceUserRefreshToken(ctx, expectedUser.ID, gomock.Any(), gomock.Any()).
			Return(nil)

		resp, err := svc.Login(ctx, payload)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Token == "" {
			t.Error("token should not be empty")
		}
		if len(resp.Bands) != 1 {
			t.Errorf("expected 1 band, got %d", len(resp.Bands))
		}
	})

	t.Run("LoginWithNoBands_OrphanUser", func(t *testing.T) {
		hashedPw, _ := hashPasswordForTest("Password123!")
		expectedUser := model.User{ID: 2, Username: "testuser", PasswordHash: hashedPw}

		mockUserRepo.EXPECT().
			FindUserByUsername(ctx, payload.Username).
			Return(expectedUser, nil)
		mockUserRepo.EXPECT().
			FindBandsByUserID(ctx, expectedUser.ID).
			Return([]model.Band{}, nil)

		mockRefreshTokenRepo.EXPECT().
			ReplaceUserRefreshToken(ctx, expectedUser.ID, gomock.Any(), gomock.Any()).
			Return(nil)

		resp, err := svc.Login(ctx, payload)

		if err != nil {
			t.Fatalf("orphan user login should succeed, got error: %v", err)
		}
		if resp.Token == "" {
			t.Error("token should not be empty")
		}
		if len(resp.Bands) != 0 {
			t.Errorf("expected 0 bands for orphan user, got %d", len(resp.Bands))
		}
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		hashedPw, _ := hashPasswordForTest("Password123!")
		expectedUser := model.User{ID: 3, Username: "testuser", PasswordHash: hashedPw}

		mockUserRepo.EXPECT().
			FindUserByUsername(ctx, payload.Username).
			Return(expectedUser, nil)

		resp, err := svc.Login(ctx, LoginPayload{Username: "testuser", Password: "wrongpassword"})

		if err == nil {
			t.Fatal("expected error for invalid password, got nil")
		}
		if resp != nil {
			t.Fatal("response should be nil on error")
		}
	})
}

func hashPasswordForTest(password string) (string, error) {
	return auth.HashPassword(password)
}

func TestUserService_CreateBand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRefreshTokenRepo := mocks.NewMockRefreshTokenRepository(ctrl)

	svc := UserService{
		UserRepo:         mockUserRepo,
		RefreshTokenRepo: mockRefreshTokenRepo,
		JWTSecret:        "testsecret",
	}

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		expectedBand := model.Band{ID: 10, Name: "New Band"}

		mockUserRepo.EXPECT().
			CreateBand(ctx, "New Band", 1).
			Return(expectedBand, nil)

		band, err := svc.CreateBand(ctx, "New Band", 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if band.ID != expectedBand.ID {
			t.Errorf("expected band ID %d, got %d", expectedBand.ID, band.ID)
		}
		if band.Name != expectedBand.Name {
			t.Errorf("expected band name %s, got %s", expectedBand.Name, band.Name)
		}
	})

	t.Run("EmptyName", func(t *testing.T) {
		band, err := svc.CreateBand(ctx, "", 1)

		if err == nil {
			t.Fatal("expected error for empty band name, got nil")
		}
		if band.ID != 0 {
			t.Errorf("expected zero Band on error, got ID=%d", band.ID)
		}
	})
}
