package model

import "time"

type Invitation struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	BandID    int       `json:"band_id"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	MaxUses   *int      `json:"max_uses"`
}
