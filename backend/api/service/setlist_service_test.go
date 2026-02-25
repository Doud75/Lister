package service

import (
	"context"
	"testing"

	"setlist/api/model"
	"setlist/api/repository/mocks"

	"go.uber.org/mock/gomock"
)

func TestSetlistService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSetlistRepository(ctrl)
	// We don't need InterludeRepo for Create
	svc := SetlistService{SetlistRepo: mockRepo}
	ctx := context.Background()
	bandID := 1

	payload := CreateSetlistPayload{
		Name:  "My Setlist",
		Color: "#FF0000",
	}

	expectedSetlist := model.Setlist{
		BandID: bandID,
		Name:   payload.Name,
		Color:  payload.Color,
	}

	// Expect GetDB to be called. We can return nil because the mock implementation of CreateSetlist
	// doesn't actually use the DB connection, it just verifies arguments.
	mockRepo.EXPECT().GetDB().Return(nil)

	// Expect CreateSetlist with the nil DB
	mockRepo.EXPECT().
		CreateSetlist(ctx, nil, payload.Name, payload.Color, bandID).
		Return(expectedSetlist, nil)

	created, err := svc.Create(ctx, payload, bandID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.Name != payload.Name {
		t.Errorf("expected name %s, got %s", payload.Name, created.Name)
	}
}

func TestSetlistService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSetlistRepository(ctrl)
	svc := SetlistService{SetlistRepo: mockRepo}
	ctx := context.Background()

	setlistID := 10
	bandID := 1
	newName := "Updated Name"
	payload := UpdateSetlistPayload{
		Name: &newName,
	}

	existingSetlist := model.Setlist{ID: setlistID, BandID: bandID, Name: "Old Name"}

	mockRepo.EXPECT().GetSetlistByID(ctx, setlistID, bandID).Return(existingSetlist, nil)

	updatedExpected := existingSetlist
	updatedExpected.Name = newName

	mockRepo.EXPECT().UpdateSetlist(ctx, updatedExpected).Return(updatedExpected, nil)

	updated, err := svc.Update(ctx, setlistID, bandID, payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != newName {
		t.Errorf("expected name %s, got %s", newName, updated.Name)
	}
}

func TestSetlistService_Duplicate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSetlistRepository(ctrl)
	// Mock Tx is generated in pgx_tx.go (package mocks)
	mockTx := mocks.NewMockTx(ctrl)

	svc := SetlistService{SetlistRepo: mockRepo}
	ctx := context.Background()

	originalID := 10
	bandID := 1
	newName := "Copy of Setlist"
	newColor := "#00FF00"

	// Expect transaction start
	mockRepo.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	// Expect Rollback (deferred)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	// Expect Get original setlist
	mockRepo.EXPECT().GetSetlistByID(ctx, originalID, bandID).Return(model.Setlist{ID: originalID, BandID: bandID}, nil)

	// Expect Get items
	songID := int32(5)
	items := []model.SetlistItem{{ID: 1, SongID: &songID}}
	mockRepo.EXPECT().GetSetlistItemsBySetlistID(ctx, originalID).Return(items, nil)

	// Expect Create new setlist within TX
	newSetlist := model.Setlist{ID: 20, Name: newName, Color: newColor}
	mockRepo.EXPECT().CreateSetlist(ctx, mockTx, newName, newColor, bandID).Return(newSetlist, nil)

	// Expect Copy items within TX
	mockRepo.EXPECT().CopyItemsToNewSetlist(ctx, mockTx, newSetlist.ID, items).Return(nil)

	// Expect Commit
	mockTx.EXPECT().Commit(ctx).Return(nil)

	result, err := svc.Duplicate(ctx, originalID, bandID, newName, newColor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != newSetlist.ID {
		t.Errorf("expected new ID %d, got %d", newSetlist.ID, result.ID)
	}
}
