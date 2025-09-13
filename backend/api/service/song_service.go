package service

import (
	"context"
	"database/sql"
	"setlist/api/model"
	"setlist/api/repository"
)

type CreateSongPayload struct {
	Title           string  `json:"title"`
	DurationSeconds *int    `json:"duration_seconds"`
	Tempo           *int    `json:"tempo"`
	SongKey         *string `json:"song_key"`
	Lyrics          *string `json:"lyrics"`
}

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

	return s.SongRepo.CreateSong(ctx, song)
}

func (s SongService) GetAllForBand(ctx context.Context, bandID int) ([]model.Song, error) {
	return s.SongRepo.GetAllSongsByBandID(ctx, bandID)
}
