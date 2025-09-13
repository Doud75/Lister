package repository

import (
	"context"
	"setlist/api/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InfoRepository struct {
	DB *pgxpool.Pool
}

func (r InfoRepository) GetBandByID(ctx context.Context, id int) (model.Band, error) {
	var band model.Band
	err := r.DB.QueryRow(ctx, `SELECT id, name FROM bands WHERE id = $1`, id).Scan(&band.ID, &band.Name)
	return band, err
}

func (r InfoRepository) GetUserByID(ctx context.Context, id int) (model.User, error) {
	var user model.User
	err := r.DB.QueryRow(ctx, `SELECT id, username FROM users WHERE id = $1`, id).Scan(&user.ID, &user.Username)
	return user, err
}
