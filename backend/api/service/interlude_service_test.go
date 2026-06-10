package service

import (
	"context"
	"errors"
	"testing"

	"setlist/api/model"
	"setlist/api/repository/mocks"

	"github.com/jackc/pgx/v5"
	"go.uber.org/mock/gomock"
)

func TestInterludeService_Create_Validation(t *testing.T) {
	svc := InterludeService{}
	ctx := context.Background()

	_, err := svc.Create(ctx, CreateInterludePayload{Title: ""}, 1)
	if !errors.Is(err, ErrInterludeTitleRequired) {
		t.Fatalf("expected ErrInterludeTitleRequired, got %v", err)
	}
}

func TestInterludeService_Update_Errors(t *testing.T) {
	ctx := context.Background()
	interludeID := 7
	bandID := 1

	t.Run("returns ErrInterludeNotFound when interlude is not in band", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockInterludeRepository(ctrl)
		svc := InterludeService{InterludeRepo: mockRepo}

		mockRepo.EXPECT().GetInterludeByID(ctx, interludeID, bandID).Return(model.Interlude{}, pgx.ErrNoRows)

		_, err := svc.Update(ctx, interludeID, bandID, UpdateInterludePayload{Title: "Title"})
		if !errors.Is(err, ErrInterludeNotFound) {
			t.Fatalf("expected ErrInterludeNotFound, got %v", err)
		}
	})

	t.Run("rejects empty title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockInterludeRepository(ctrl)
		svc := InterludeService{InterludeRepo: mockRepo}

		mockRepo.EXPECT().GetInterludeByID(ctx, interludeID, bandID).Return(model.Interlude{ID: interludeID}, nil)

		_, err := svc.Update(ctx, interludeID, bandID, UpdateInterludePayload{Title: ""})
		if !errors.Is(err, ErrInterludeTitleRequired) {
			t.Fatalf("expected ErrInterludeTitleRequired, got %v", err)
		}
	})

	t.Run("propagates unexpected repository errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockInterludeRepository(ctrl)
		svc := InterludeService{InterludeRepo: mockRepo}

		dbErr := errors.New("connection lost")
		mockRepo.EXPECT().GetInterludeByID(ctx, interludeID, bandID).Return(model.Interlude{}, dbErr)

		_, err := svc.Update(ctx, interludeID, bandID, UpdateInterludePayload{Title: "Title"})
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected raw db error, got %v", err)
		}
	})
}
