package service

import (
	"context"
	"encoding/json"
	"setlist/api/model"
	"setlist/api/repository"
	"setlist/cache"
	"time"

	"github.com/redis/go-redis/v9"
)

const songCacheTTL = time.Hour

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
	Cache    *redis.Client
}

func ptrInt32(v *int) *int32 {
	if v == nil {
		return nil
	}
	v32 := int32(*v)
	return &v32
}

func (s SongService) buildSong(bandID int, payload CreateSongPayload) model.Song {
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
	return song
}

func (s SongService) Create(ctx context.Context, payload CreateSongPayload, bandID int) (model.Song, error) {
	song := s.buildSong(bandID, payload)

	created, err := s.SongRepo.CreateSong(ctx, song)
	if err != nil {
		return model.Song{}, err
	}

	cache.Delete(ctx, s.Cache, cache.SongKey(bandID))

	return created, nil
}

func (s SongService) GetAllForBand(ctx context.Context, bandID int) ([]model.Song, error) {
	key := cache.SongKey(bandID)

	if data, ok := cache.Get(ctx, s.Cache, key); ok {
		var songs []model.Song
		if err := json.Unmarshal([]byte(data), &songs); err == nil {
			return songs, nil
		}
	}

	songs, err := s.SongRepo.GetAllSongsByBandID(ctx, bandID)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(songs); err == nil {
		cache.Set(ctx, s.Cache, key, string(data), songCacheTTL)
	}

	return songs, nil
}

func (s SongService) GetByID(ctx context.Context, id int, bandID int) (model.Song, error) {
	return s.SongRepo.GetSongByID(ctx, id, bandID)
}

func (s SongService) Update(ctx context.Context, id int, bandID int, payload UpdateSongPayload) (model.Song, error) {
	song := s.buildSong(bandID, payload)
	song.ID = id

	updated, err := s.SongRepo.UpdateSong(ctx, song)
	if err != nil {
		return model.Song{}, err
	}

	cache.Delete(ctx, s.Cache, cache.SongKey(bandID))

	return updated, nil
}

func (s SongService) SoftDelete(ctx context.Context, id int, bandID int) error {
	if err := s.SongRepo.SoftDeleteSong(ctx, id, bandID); err != nil {
		return err
	}

	cache.Delete(ctx, s.Cache, cache.SongKey(bandID))

	return nil
}
