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

func (r *UserRepository) CreateBandAndUser(ctx context.Context, bandName, username, passwordHash string) (model.User, model.Band, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return model.User{}, model.Band{}, err
	}
	defer tx.Rollback(ctx)

	var band model.Band
	bandQuery := `INSERT INTO bands (name) VALUES ($1) RETURNING id, name`
	err = tx.QueryRow(ctx, bandQuery, bandName).Scan(&band.ID, &band.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.User{}, model.Band{}, ErrDuplicateBand
		}
		return model.User{}, model.Band{}, err
	}

	var user model.User
	userQuery := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, created_at`
	err = tx.QueryRow(ctx, userQuery, username, passwordHash).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.User{}, model.Band{}, ErrDuplicateUsername
		}
		return model.User{}, model.Band{}, err
	}

	linkQuery := `INSERT INTO band_users (user_id, band_id, role) VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, linkQuery, user.ID, band.ID, "admin")
	if err != nil {
		return model.User{}, model.Band{}, err
	}

	return user, band, tx.Commit(ctx)
}

func (r *UserRepository) CreateUserForExistingBand(ctx context.Context, bandName, username, passwordHash string) (model.User, model.Band, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return model.User{}, model.Band{}, err
	}
	defer tx.Rollback(ctx)

	var band model.Band
	bandQuery := `SELECT id, name FROM bands WHERE name = $1`
	err = tx.QueryRow(ctx, bandQuery, bandName).Scan(&band.ID, &band.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.Band{}, errors.New("band not found")
		}
		return model.User{}, model.Band{}, err
	}

	var user model.User
	userQuery := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, created_at`
	err = tx.QueryRow(ctx, userQuery, username, passwordHash).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.User{}, model.Band{}, ErrDuplicateUsername
		}
		return model.User{}, model.Band{}, err
	}

	linkQuery := `INSERT INTO band_users (user_id, band_id, role) VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, linkQuery, user.ID, band.ID, "member")
	if err != nil {
		return model.User{}, model.Band{}, err
	}

	return user, band, tx.Commit(ctx)
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	query := `SELECT id, password_hash, username FROM users WHERE username = $1`
	err := r.DB.QueryRow(ctx, query, username).Scan(&user.ID, &user.PasswordHash, &user.Username)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) FindBandsByUserID(ctx context.Context, userID int) ([]model.Band, error) {
	var bands []model.Band
	query := `
		SELECT b.id, b.name FROM bands b
		JOIN band_users bu ON b.id = bu.band_id
		WHERE bu.user_id = $1
		ORDER BY b.name
	`
	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var band model.Band
		if err := rows.Scan(&band.ID, &band.Name); err != nil {
			return nil, err
		}
		bands = append(bands, band)
	}

	return bands, rows.Err()
}

func (r *UserRepository) IsUserInBand(ctx context.Context, userID int, bandID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM band_users WHERE user_id = $1 AND band_id = $2)`
	err := r.DB.QueryRow(ctx, query, userID, bandID).Scan(&exists)
	return exists, err
}