package model

type SetlistItem struct {
	ID                        int     `json:"id"`
	SetlistID                 int     `json:"setlist_id"`
	Position                  int     `json:"position"`
	ItemType                  string  `json:"item_type"`
	SongID                    *int32  `json:"song_id,omitempty"`
	InterludeID               *int32  `json:"interlude_id,omitempty"`
	Notes                     *string `json:"notes"`
	TransitionDurationSeconds int     `json:"transition_duration_seconds"`
	Title                     *string `json:"title,omitempty"`
	DurationSeconds           *int32  `json:"duration_seconds,omitempty"`
	Tempo                     *int32  `json:"tempo,omitempty"`
	Speaker                   *string `json:"speaker,omitempty"`
	Script                    *string `json:"script,omitempty"`
	SongKey                   *string `json:"song_key,omitempty"`
	Links                     *string `json:"links,omitempty"`
}
