package repository

import (
	"context"
	"errors"
	"setlist/api/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDuplicateBand     = errors.New("band with this name already exists")
	ErrDuplicateUsername = errors.New("username already exists")
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func (r *UserRepository) CreateBandAndUser(ctx context.Context, bandName, username, passwordHash string) (model.User, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return model.User{}, err
	}
	defer tx.Rollback(ctx)

	var bandID int
	bandQuery := `INSERT INTO bands (name) VALUES ($1) RETURNING id`
	err = tx.QueryRow(ctx, bandQuery, bandName).Scan(&bandID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.User{}, ErrDuplicateBand
		}
		return model.User{}, err
	}

	var user model.User
	userQuery := `INSERT INTO users (band_id, username, password_hash) VALUES ($1, $2, $3) RETURNING id, band_id, username, created_at`
	err = tx.QueryRow(ctx, userQuery, bandID, username, passwordHash).Scan(&user.ID, &user.BandID, &user.Username, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.User{}, ErrDuplicateUsername
		}
		return model.User{}, err
	}

	return user, tx.Commit(ctx)
}

func (r *UserRepository) CreateUserForExistingBand(ctx context.Context, bandName, username, passwordHash string) (model.User, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return model.User{}, err
	}
	defer tx.Rollback(ctx)

	var band model.Band
	bandQuery := `SELECT id, name FROM bands WHERE name = $1`
	err = tx.QueryRow(ctx, bandQuery, bandName).Scan(&band.ID, &band.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errors.New("band not found")
		}
		return model.User{}, err
	}

	var user model.User
	userQuery := `INSERT INTO users (band_id, username, password_hash) VALUES ($1, $2, $3) RETURNING id, band_id, username, created_at`
	err = tx.QueryRow(ctx, userQuery, band.ID, username, passwordHash).Scan(&user.ID, &user.BandID, &user.Username, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.User{}, ErrDuplicateUsername
		}
		return model.User{}, err
	}

	return user, tx.Commit(ctx)
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	query := `SELECT id, band_id, password_hash, username FROM users WHERE username = $1`
	err := r.DB.QueryRow(ctx, query, username).Scan(&user.ID, &user.BandID, &user.PasswordHash, &user.Username)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
