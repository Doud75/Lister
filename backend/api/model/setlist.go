package model

import "time"

type Setlist struct {
	ID        int       `json:"id"`
	BandID    int       `json:"band_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}
