package model

import (
	"database/sql"
)

type SetlistItem struct {
	ID                        int            `json:"id"`
	SetlistID                 int            `json:"setlist_id"`
	Position                  int            `json:"position"`
	ItemType                  string         `json:"item_type"`
	SongID                    sql.NullInt32  `json:"song_id,omitempty"`
	InterludeID               sql.NullInt32  `json:"interlude_id,omitempty"`
	Notes                     sql.NullString `json:"notes"`
	TransitionDurationSeconds int            `json:"transition_duration_seconds"`
	Title                     sql.NullString `json:"title,omitempty"`
	DurationSeconds           sql.NullInt32  `json:"duration_seconds,omitempty"`
	Tempo                     sql.NullInt32  `json:"tempo,omitempty"`
	Speaker                   sql.NullString `json:"speaker,omitempty"`
	Script                    sql.NullString `json:"script,omitempty"`
	SongKey                   sql.NullString `json:"song_key,omitempty"`
	Links                     sql.NullString `json:"links,omitempty"`
}
