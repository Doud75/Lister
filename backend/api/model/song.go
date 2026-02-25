package model

import (
	"encoding/json"
	"time"
)

type Song struct {
	ID              int             `json:"id"`
	BandID          int             `json:"band_id"`
	Title           string          `json:"title"`
	DurationSeconds *int32          `json:"duration_seconds"`
	Tempo           *int32          `json:"tempo"`
	SongKey         *string         `json:"song_key"`
	Lyrics          *string         `json:"lyrics"`
	Chords          *string         `json:"chords"`
	AlbumName       *string         `json:"album_name"`
	Instrumentation json.RawMessage `json:"instrumentation"`
	Notes           *string         `json:"notes"`
	Links           *string         `json:"links"`
	CreatedAt       time.Time       `json:"created_at"`
}
