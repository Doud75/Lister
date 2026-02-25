package service

import (
	"context"
	"encoding/json"
	"setlist/api/model"
	"setlist/api/repository"
)

type CreateSongPayload struct {
	Title           string           `json:"title"`
	DurationSeconds *int             `json:"duration_seconds"`
	Tempo           *int             `json:"tempo"`
	SongKey         *string          `json:"song_key"`
	Lyrics          *string          `json:"lyrics"`
	AlbumName       *string          `json:"album_name"`
	Instrumentation *json.RawMessage `json:"instrumentation"`
	Links           *string          `json:"links"`
}

type UpdateSongPayload = CreateSongPayload

type SongService struct {
	SongRepo repository.SongRepository
}

func ptrInt32(v *int) *int32 {
	if v == nil {
		return nil
	}
	i := int32(*v)
	return &i
}

func (s SongService) Create(ctx context.Context, payload CreateSongPayload, bandID int) (model.Song, error) {
	song := model.Song{
		BandID:          bandID,
		Title:           payload.Title,
		DurationSeconds: ptrInt32(payload.DurationSeconds),
		Tempo:           ptrInt32(payload.Tempo),
		SongKey:         payload.SongKey,
		Lyrics:          payload.Lyrics,
		AlbumName:       payload.AlbumName,
		Links:           payload.Links,
	}

	if payload.Instrumentation != nil {
		song.Instrumentation = *payload.Instrumentation
	} else {
		song.Instrumentation = json.RawMessage("null")
	}

	return s.SongRepo.CreateSong(ctx, song)
}

func (s SongService) GetAllForBand(ctx context.Context, bandID int) ([]model.Song, error) {
	return s.SongRepo.GetAllSongsByBandID(ctx, bandID)
}

func (s SongService) GetByID(ctx context.Context, id int, bandID int) (model.Song, error) {
	return s.SongRepo.GetSongByID(ctx, id, bandID)
}

func (s SongService) Update(ctx context.Context, id int, bandID int, payload UpdateSongPayload) (model.Song, error) {
	song := model.Song{
		ID:              id,
		BandID:          bandID,
		Title:           payload.Title,
		DurationSeconds: ptrInt32(payload.DurationSeconds),
		Tempo:           ptrInt32(payload.Tempo),
		SongKey:         payload.SongKey,
		Lyrics:          payload.Lyrics,
		AlbumName:       payload.AlbumName,
		Links:           payload.Links,
	}

	if payload.Instrumentation != nil {
		song.Instrumentation = *payload.Instrumentation
	} else {
		song.Instrumentation = json.RawMessage("null")
	}

	return s.SongRepo.UpdateSong(ctx, song)
}

func (s SongService) SoftDelete(ctx context.Context, id int, bandID int) error {
	return s.SongRepo.SoftDeleteSong(ctx, id, bandID)
}
