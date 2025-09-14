package repository

import (
	"context"
	"setlist/api/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InterludeRepository struct {
	DB *pgxpool.Pool
}

func (r InterludeRepository) CreateInterlude(ctx context.Context, interlude model.Interlude) (model.Interlude, error) {
	query := `
		INSERT INTO interludes (band_id, title, speaker, script, duration_seconds)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := r.DB.QueryRow(ctx, query,
		interlude.BandID,
		interlude.Title,
		interlude.Speaker,
		interlude.Script,
		interlude.DurationSeconds,
	).Scan(&interlude.ID, &interlude.CreatedAt)

	return interlude, err
}

func (r InterludeRepository) GetAllInterludesByBandID(ctx context.Context, bandID int) ([]model.Interlude, error) {
	interludes := make([]model.Interlude, 0)
	query := `SELECT id, title FROM interludes WHERE band_id = $1 ORDER BY title ASC`

	rows, err := r.DB.Query(ctx, query, bandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var interlude model.Interlude
		if err := rows.Scan(&interlude.ID, &interlude.Title); err != nil {
			return nil, err
		}
		interludes = append(interludes, interlude)
	}
	return interludes, rows.Err()
}

func (r InterludeRepository) GetInterludeByID(ctx context.Context, id int, bandID int) (model.Interlude, error) {
	var interlude model.Interlude
	query := `SELECT id, band_id, title, speaker, script, duration_seconds FROM interludes WHERE id = $1 AND band_id = $2`
	err := r.DB.QueryRow(ctx, query, id, bandID).Scan(&interlude.ID, &interlude.BandID, &interlude.Title, &interlude.Speaker, &interlude.Script, &interlude.DurationSeconds)
	return interlude, err
}

func (r InterludeRepository) UpdateInterlude(ctx context.Context, interlude model.Interlude) error {
	query := `UPDATE interludes SET title = $1, speaker = $2, script = $3, duration_seconds = $4 WHERE id = $5 AND band_id = $6`
	_, err := r.DB.Exec(ctx, query, interlude.Title, interlude.Speaker, interlude.Script, interlude.DurationSeconds, interlude.ID, interlude.BandID)
	return err
}

func (r InterludeRepository) DeleteInterlude(ctx context.Context, id int, bandID int) error {
	query := `DELETE FROM interludes WHERE id = $1 AND band_id = $2`
	_, err := r.DB.Exec(ctx, query, id, bandID)
	return err
}
