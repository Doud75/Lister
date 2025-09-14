package service

import (
	"context"
	"database/sql"
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

func (s SongService) Create(ctx context.Context, payload CreateSongPayload, bandID int) (model.Song, error) {
	song := model.Song{
		BandID: bandID,
		Title:  payload.Title,
	}

	if payload.DurationSeconds != nil {
		song.DurationSeconds = sql.NullInt32{Int32: int32(*payload.DurationSeconds), Valid: true}
	}
	if payload.Tempo != nil {
		song.Tempo = sql.NullInt32{Int32: int32(*payload.Tempo), Valid: true}
	}
	if payload.SongKey != nil {
		song.SongKey = sql.NullString{String: *payload.SongKey, Valid: true}
	}
	if payload.Lyrics != nil {
		song.Lyrics = sql.NullString{String: *payload.Lyrics, Valid: true}
	}
	if payload.AlbumName != nil {
		song.AlbumName = sql.NullString{String: *payload.AlbumName, Valid: true}
	}
	if payload.Instrumentation != nil {
		song.Instrumentation = *payload.Instrumentation
	} else {
		song.Instrumentation = json.RawMessage("null")
	}
	if payload.Links != nil {
		song.Links = sql.NullString{String: *payload.Links, Valid: true}
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
		ID:     id,
		BandID: bandID,
		Title:  payload.Title,
	}

	if payload.DurationSeconds != nil {
		song.DurationSeconds = sql.NullInt32{Int32: int32(*payload.DurationSeconds), Valid: true}
	}
	if payload.Tempo != nil {
		song.Tempo = sql.NullInt32{Int32: int32(*payload.Tempo), Valid: true}
	}
	if payload.SongKey != nil {
		song.SongKey = sql.NullString{String: *payload.SongKey, Valid: true}
	}
	if payload.Lyrics != nil {
		song.Lyrics = sql.NullString{String: *payload.Lyrics, Valid: true}
	}
	if payload.AlbumName != nil {
		song.AlbumName = sql.NullString{String: *payload.AlbumName, Valid: true}
	}
	if payload.Instrumentation != nil {
		song.Instrumentation = *payload.Instrumentation
	} else {
		song.Instrumentation = json.RawMessage("null")
	}
	if payload.Links != nil {
		song.Links = sql.NullString{String: *payload.Links, Valid: true}
	}

	return s.SongRepo.UpdateSong(ctx, song)
}

func (s SongService) SoftDelete(ctx context.Context, id int, bandID int) error {
	return s.SongRepo.SoftDeleteSong(ctx, id, bandID)
}
