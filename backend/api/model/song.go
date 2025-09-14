package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Song struct {
	ID              int             `json:"id"`
	BandID          int             `json:"band_id"`
	Title           string          `json:"title"`
	DurationSeconds sql.NullInt32   `json:"duration_seconds"`
	Tempo           sql.NullInt32   `json:"tempo"`
	SongKey         sql.NullString  `json:"song_key"`
	Lyrics          sql.NullString  `json:"lyrics"`
	Chords          sql.NullString  `json:"chords"`
	AlbumName       sql.NullString  `json:"album_name"`
	Instrumentation json.RawMessage `json:"instrumentation"`
	Notes           sql.NullString  `json:"notes"`
	Links           sql.NullString  `json:"links"`
	CreatedAt       time.Time       `json:"created_at"`
}
