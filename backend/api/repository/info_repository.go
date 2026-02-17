package repository

import (
	"context"
	"setlist/api/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InfoRepository interface {
	GetBandByID(ctx context.Context, id int) (model.Band, error)
	GetUserByID(ctx context.Context, id int) (model.User, error)
}

type PgInfoRepository struct {
	DB *pgxpool.Pool
}

func (r PgInfoRepository) GetBandByID(ctx context.Context, id int) (model.Band, error) {
	var band model.Band
	err := r.DB.QueryRow(ctx, `SELECT id, name FROM bands WHERE id = $1`, id).Scan(&band.ID, &band.Name)
	return band, err
}

func (r PgInfoRepository) GetUserByID(ctx context.Context, id int) (model.User, error) {
	var user model.User
	err := r.DB.QueryRow(ctx, `SELECT id, username FROM users WHERE id = $1`, id).Scan(&user.ID, &user.Username)
	return user, err
}
