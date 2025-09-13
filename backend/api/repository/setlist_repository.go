package repository

import (
	"context"
	"database/sql"
	"setlist/api/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SetlistRepository struct {
	DB *pgxpool.Pool
}

func (r SetlistRepository) CreateSetlist(ctx context.Context, name, color string, bandID int) (model.Setlist, error) {
	var setlist model.Setlist
	query := `
		INSERT INTO setlists (name, color, band_id)
		VALUES ($1, $2, $3)
		RETURNING id, band_id, name, color, created_at
	`
	err := r.DB.QueryRow(ctx, query, name, color, bandID).Scan(
		&setlist.ID, &setlist.BandID, &setlist.Name, &setlist.Color, &setlist.CreatedAt,
	)
	return setlist, err
}

func (r SetlistRepository) GetSetlistsByBandID(ctx context.Context, bandID int) ([]model.Setlist, error) {
	setlists := make([]model.Setlist, 0)
	query := `
		SELECT id, band_id, name, color, created_at
		FROM setlists
		WHERE band_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.DB.Query(ctx, query, bandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var setlist model.Setlist
		if err := rows.Scan(&setlist.ID, &setlist.BandID, &setlist.Name, &setlist.Color, &setlist.CreatedAt); err != nil {
			return setlists, err
		}
		setlists = append(setlists, setlist)
	}

	if err := rows.Err(); err != nil {
		return setlists, err
	}

	return setlists, nil
}

func (r SetlistRepository) GetSetlistByID(ctx context.Context, id int, bandID int) (model.Setlist, error) {
	var setlist model.Setlist
	query := `SELECT id, band_id, name, color, created_at FROM setlists WHERE id = $1 AND band_id = $2`
	err := r.DB.QueryRow(ctx, query, id, bandID).Scan(&setlist.ID, &setlist.BandID, &setlist.Name, &setlist.Color, &setlist.CreatedAt)
	return setlist, err
}

func (r SetlistRepository) GetSetlistItemsBySetlistID(ctx context.Context, setlistID int) ([]model.SetlistItem, error) {
	items := make([]model.SetlistItem, 0)
	query := `
		SELECT
			si.id, si.setlist_id, si.position, si.item_type,
			si.song_id, si.interlude_id, si.notes, si.transition_duration_seconds,
			COALESCE(s.title, i.title) as title,
			COALESCE(s.duration_seconds, i.duration_seconds) as duration_seconds,
			s.tempo
		FROM setlist_items si
		LEFT JOIN songs s ON si.song_id = s.id
		LEFT JOIN interludes i ON si.interlude_id = i.id
		WHERE si.setlist_id = $1
		ORDER BY si.position ASC
	`
	rows, err := r.DB.Query(ctx, query, setlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.SetlistItem
		err := rows.Scan(
			&item.ID, &item.SetlistID, &item.Position, &item.ItemType,
			&item.SongID, &item.InterludeID, &item.Notes, &item.TransitionDurationSeconds,
			&item.Title, &item.DurationSeconds, &item.Tempo, // <- Ces champs viennent du modèle corrigé
		)
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return items, err
	}
	return items, nil
}

func (r SetlistRepository) AddItemToSetlist(ctx context.Context, item model.SetlistItem) (model.SetlistItem, error) {
	var maxPosition sql.NullInt32
	posQuery := `SELECT MAX(position) FROM setlist_items WHERE setlist_id = $1`
	r.DB.QueryRow(ctx, posQuery, item.SetlistID).Scan(&maxPosition)

	nextPos := 0
	if maxPosition.Valid {
		nextPos = int(maxPosition.Int32) + 1
	}
	item.Position = nextPos

	insertQuery := `INSERT INTO setlist_items (setlist_id, position, item_type, song_id, interlude_id, notes, transition_duration_seconds)
					VALUES ($1, $2, $3, $4, $5, $6, $7)
					RETURNING id`

	err := r.DB.QueryRow(ctx, insertQuery,
		item.SetlistID,
		item.Position,
		item.ItemType,
		item.SongID,
		item.InterludeID,
		item.Notes,
		item.TransitionDurationSeconds,
	).Scan(&item.ID)

	return item, err
}

func (r SetlistRepository) UpdateItemOrder(ctx context.Context, setlistID int, itemIDs []int) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "SET CONSTRAINTS unique_position_in_setlist DEFERRED"); err != nil {
		return err
	}

	query := `UPDATE setlist_items SET position = $1 WHERE id = $2 AND setlist_id = $3`

	for i, id := range itemIDs {
		if _, err := tx.Exec(ctx, query, i, id, setlistID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
