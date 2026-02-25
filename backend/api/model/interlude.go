package model

import (
	"time"
)

type Interlude struct {
	ID              int       `json:"id"`
	BandID          int       `json:"band_id"`
	Title           string    `json:"title"`
	Speaker         *string   `json:"speaker"`
	Script          *string   `json:"script"`
	DurationSeconds *int32    `json:"duration_seconds"`
	CreatedAt       time.Time `json:"created_at"`
}
