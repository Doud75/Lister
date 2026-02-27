package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"setlist/api/middleware"
	"setlist/api/model"
	"setlist/api/repository/mocks"
	"setlist/api/service"

	"go.uber.org/mock/gomock"
)

func TestInvitationHandler_CreateInvitation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvRepo := mocks.NewMockInvitationRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	invSvc := service.InvitationService{
		InvitationRepo: mockInvRepo,
		UserRepo:       mockUserRepo,
	}
	h := InvitationHandler{InvitationService: invSvc}

	t.Run("Success", func(t *testing.T) {
		payload := CreateInvitationPayload{Role: "member"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/api/bands/1/invitations", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.BandIDKey, 1)
		req = req.WithContext(ctx)

		mockInvRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(nil)

		w := httptest.NewRecorder()
		err := h.CreateInvitation(w, req)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		var resp CreateInvitationResponse
		json.NewDecoder(w.Body).Decode(&resp)
		if len(resp.Token) == 0 {
			t.Error("expected a token")
		}
	})
}

func TestInvitationHandler_GetInvitation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvRepo := mocks.NewMockInvitationRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	invSvc := service.InvitationService{
		InvitationRepo: mockInvRepo,
		UserRepo:       mockUserRepo,
	}
	h := InvitationHandler{InvitationService: invSvc}

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/invitations/testtoken", nil)
		req.SetPathValue("token", "testtoken")

		detail := &model.InvitationDetail{
			Token:     "testtoken",
			Role:      "member",
			BandName:  "Test Band",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		mockInvRepo.EXPECT().
			GetByToken(gomock.Any(), "testtoken").
			Return(detail, nil)

		w := httptest.NewRecorder()
		err := h.GetInvitation(w, req)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp GetInvitationResponse
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.BandName != "Test Band" {
			t.Errorf("expected band Test Band, got %s", resp.BandName)
		}
	})

	t.Run("Expired", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/invitations/testtoken", nil)
		req.SetPathValue("token", "testtoken")

		detail := &model.InvitationDetail{
			Token:     "testtoken",
			Role:      "member",
			BandName:  "Test Band",
			ExpiresAt: time.Now().Add(-24 * time.Hour),
		}

		mockInvRepo.EXPECT().
			GetByToken(gomock.Any(), "testtoken").
			Return(detail, nil)

		w := httptest.NewRecorder()
		err := h.GetInvitation(w, req)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestInvitationHandler_AcceptInvitation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvRepo := mocks.NewMockInvitationRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	invSvc := service.InvitationService{
		InvitationRepo: mockInvRepo,
		UserRepo:       mockUserRepo,
	}
	h := InvitationHandler{InvitationService: invSvc}

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/invitations/testtoken/accept", nil)
		req.SetPathValue("token", "testtoken")
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, 42)
		req = req.WithContext(ctx)

		detail := &model.InvitationDetail{
			Token:     "testtoken",
			BandID:    10,
			Role:      "member",
			BandName:  "Test Band",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		mockInvRepo.EXPECT().
			GetByToken(gomock.Any(), "testtoken").
			Return(detail, nil)

		mockUserRepo.EXPECT().
			IsUserInBand(gomock.Any(), 42, 10).
			Return(false, nil)

		mockUserRepo.EXPECT().
			AddUserToBand(gomock.Any(), 42, 10, "member").
			Return(nil)

		w := httptest.NewRecorder()
		err := h.AcceptInvitation(w, req)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp AcceptInvitationResponse
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.BandID != 10 {
			t.Errorf("expected band ID 10, got %d", resp.BandID)
		}
	})
}
