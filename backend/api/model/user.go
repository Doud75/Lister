package model

import "time"

type User struct {
	ID           int       `json:"id"`
	BandID       int       `json:"band_id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
