package model

import (
	"database/sql"
	"time"
)

type Interlude struct {
	ID              int            `json:"id"`
	BandID          int            `json:"band_id"`
	Title           string         `json:"title"`
	Speaker         sql.NullString `json:"speaker"`
	Script          sql.NullString `json:"script"`
	DurationSeconds sql.NullInt32  `json:"duration_seconds"`
	CreatedAt       time.Time      `json:"created_at"`
}
