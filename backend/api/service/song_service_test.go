package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"setlist/api/model"
	"setlist/api/repository/mocks"

	"go.uber.org/mock/gomock"
)

func ptr32(v int32) *int32    { return &v }
func ptrStr(v string) *string { return &v }

func TestSongService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSongRepository(ctrl)
	svc := SongService{SongRepo: mockRepo}
	ctx := context.Background()
	bandID := 1

	duration := 180
	tempo := 120
	key := "Am"
	lyrics := "La la la"
	payload := CreateSongPayload{
		Title:           "New Song",
		DurationSeconds: &duration,
		Tempo:           &tempo,
		SongKey:         &key,
		Lyrics:          &lyrics,
	}

	expectedSong := model.Song{
		BandID:          bandID,
		Title:           payload.Title,
		DurationSeconds: ptr32(int32(duration)),
		Tempo:           ptr32(int32(tempo)),
		SongKey:         ptrStr(key),
		Lyrics:          ptrStr(lyrics),
		Instrumentation: json.RawMessage("null"),
	}

	mockRepo.EXPECT().
		CreateSong(ctx, expectedSong).
		Return(model.Song{ID: 10, Title: "New Song"}, nil)

	created, err := svc.Create(ctx, payload, bandID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID != 10 {
		t.Errorf("expected ID 10, got %d", created.ID)
	}
}

func TestSongService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSongRepository(ctrl)
	svc := SongService{SongRepo: mockRepo}
	ctx := context.Background()

	songID := 10
	bandID := 1
	expectedSong := model.Song{ID: songID, BandID: bandID, Title: "My Song"}

	t.Run("Found", func(t *testing.T) {
		mockRepo.EXPECT().GetSongByID(ctx, songID, bandID).Return(expectedSong, nil)

		song, err := svc.GetByID(ctx, songID, bandID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if song.Title != expectedSong.Title {
			t.Errorf("expected title %s, got %s", expectedSong.Title, song.Title)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetSongByID(ctx, songID, bandID).Return(model.Song{}, errors.New("not found"))

		_, err := svc.GetByID(ctx, songID, bandID)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestSongService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSongRepository(ctrl)
	svc := SongService{SongRepo: mockRepo}
	ctx := context.Background()

	songID := 10
	bandID := 1
	newTitle := "Updated Title"
	payload := UpdateSongPayload{
		Title: newTitle,
	}

	expectedSong := model.Song{
		ID:              songID,
		BandID:          bandID,
		Title:           newTitle,
		Instrumentation: json.RawMessage("null"), // Default null if not provided
	}

	mockRepo.EXPECT().UpdateSong(ctx, expectedSong).Return(expectedSong, nil)

	updated, err := svc.Update(ctx, songID, bandID, payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Title != newTitle {
		t.Errorf("expected title %s, got %s", newTitle, updated.Title)
	}
}

func TestSongService_SoftDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSongRepository(ctrl)
	svc := SongService{SongRepo: mockRepo}
	ctx := context.Background()

	songID := 10
	bandID := 1

	mockRepo.EXPECT().SoftDeleteSong(ctx, songID, bandID).Return(nil)

	err := svc.SoftDelete(ctx, songID, bandID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSongService_GetAllForBand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSongRepository(ctrl)
	svc := SongService{SongRepo: mockRepo}
	ctx := context.Background()
	bandID := 1

	expectedSongs := []model.Song{
		{ID: 1, Title: "Song 1"},
		{ID: 2, Title: "Song 2"},
	}

	mockRepo.EXPECT().GetAllSongsByBandID(ctx, bandID).Return(expectedSongs, nil)

	songs, err := svc.GetAllForBand(ctx, bandID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(songs) != 2 {
		t.Errorf("expected 2 songs, got %d", len(songs))
	}
}
