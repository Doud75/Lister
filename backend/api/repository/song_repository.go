package repository

import (
	"context"
	"setlist/api/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SongRepository struct {
	DB *pgxpool.Pool
}

func (r SongRepository) CreateSong(ctx context.Context, song model.Song) (model.Song, error) {
	query := `
		INSERT INTO songs (band_id, title, duration_seconds, tempo, song_key, lyrics)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	err := r.DB.QueryRow(ctx, query,
		song.BandID,
		song.Title,
		song.DurationSeconds,
		song.Tempo,
		song.SongKey,
		song.Lyrics,
	).Scan(&song.ID, &song.CreatedAt)

	return song, err
}

func (r SongRepository) GetAllSongsByBandID(ctx context.Context, bandID int) ([]model.Song, error) {
	songs := make([]model.Song, 0)
	query := `SELECT id, title FROM songs WHERE band_id = $1 ORDER BY title ASC`

	rows, err := r.DB.Query(ctx, query, bandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song model.Song
		if err := rows.Scan(&song.ID, &song.Title); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, rows.Err()
}

func (r SongRepository) GetSongByID(ctx context.Context, id int, bandID int) (model.Song, error) {
	var song model.Song
	query := `SELECT id, band_id, title, duration_seconds, tempo, song_key FROM songs WHERE id = $1 AND band_id = $2`
	err := r.DB.QueryRow(ctx, query, id, bandID).Scan(&song.ID, &song.BandID, &song.Title, &song.DurationSeconds, &song.Tempo, &song.SongKey)
	return song, err
}

func (r SongRepository) UpdateSong(ctx context.Context, song model.Song) error {
	query := `UPDATE songs SET title = $1, duration_seconds = $2, tempo = $3, song_key = $4 WHERE id = $5 AND band_id = $6`
	_, err := r.DB.Exec(ctx, query, song.Title, song.DurationSeconds, song.Tempo, song.SongKey, song.ID, song.BandID)
	return err
}

func (r SongRepository) DeleteSong(ctx context.Context, id int, bandID int) error {
	query := `DELETE FROM songs WHERE id = $1 AND band_id = $2`
	_, err := r.DB.Exec(ctx, query, id, bandID)
	return err
}
