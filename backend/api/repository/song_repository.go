package repository

import (
	"context"
	"database/sql"
	"setlist/api/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SongRepository struct {
	DB *pgxpool.Pool
}

func (r SongRepository) CreateSong(ctx context.Context, song model.Song) (model.Song, error) {
	query := `
		INSERT INTO songs (
			band_id, title, duration_seconds, tempo, song_key, lyrics, album_name, instrumentation, links
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`
	err := r.DB.QueryRow(ctx, query,
		song.BandID,
		song.Title,
		song.DurationSeconds,
		song.Tempo,
		song.SongKey,
		song.Lyrics,
		song.AlbumName,
		song.Instrumentation,
		song.Links,
	).Scan(&song.ID, &song.CreatedAt)

	return song, err
}

func (r SongRepository) GetAllSongsByBandID(ctx context.Context, bandID int) ([]model.Song, error) {
	songs := make([]model.Song, 0)
	query := `SELECT id, title, album_name FROM songs WHERE band_id = $1 AND is_deleted = FALSE ORDER BY album_name ASC, title ASC`

	rows, err := r.DB.Query(ctx, query, bandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song model.Song
		if err := rows.Scan(&song.ID, &song.Title, &song.AlbumName); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, rows.Err()
}

func (r SongRepository) GetSongByID(ctx context.Context, id int, bandID int) (model.Song, error) {
	var song model.Song
	query := `
		SELECT 
			id, band_id, title, duration_seconds, tempo, song_key, lyrics, album_name, instrumentation, links, created_at
		FROM songs 
		WHERE id = $1 AND band_id = $2 AND is_deleted = FALSE
	`
	err := r.DB.QueryRow(ctx, query, id, bandID).Scan(
		&song.ID, &song.BandID, &song.Title, &song.DurationSeconds, &song.Tempo, &song.SongKey, &song.Lyrics,
		&song.AlbumName, &song.Instrumentation, &song.Links, &song.CreatedAt,
	)
	return song, err
}

func (r SongRepository) UpdateSong(ctx context.Context, song model.Song) (model.Song, error) {
	query := `
		UPDATE songs SET
			title = $1, duration_seconds = $2, tempo = $3, song_key = $4, lyrics = $5,
			album_name = $6, instrumentation = $7, links = $8, updated_at = NOW()
		WHERE id = $9 AND band_id = $10
		RETURNING id, created_at, updated_at
	`
	err := r.DB.QueryRow(ctx, query,
		song.Title, song.DurationSeconds, song.Tempo, song.SongKey, song.Lyrics,
		song.AlbumName, song.Instrumentation, song.Links,
		song.ID, song.BandID,
	).Scan(&song.ID, &song.CreatedAt, &song.CreatedAt)

	return song, err
}

func (r SongRepository) SoftDeleteSong(ctx context.Context, id int, bandID int) error {
	query := `UPDATE songs SET is_deleted = TRUE WHERE id = $1 AND band_id = $2`
	cmdTag, err := r.DB.Exec(ctx, query, id, bandID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}
