package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	DB *pgxpool.Pool
}

func (r *RefreshTokenRepository) StoreRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (r *RefreshTokenRepository) FindRefreshToken(ctx context.Context, tokenHash string) (userID int, expiresAt time.Time, err error) {
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash = $1`
	err = r.DB.QueryRow(ctx, query, tokenHash).Scan(&userID, &expiresAt)
	return
}

func (r *RefreshTokenRepository) GetUserTokenHashes(ctx context.Context, userID int) ([]string, error) {
	query := `SELECT token_hash FROM refresh_tokens WHERE user_id = $1`
	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hashes []string
	for rows.Next() {
		var hash string
		if err := rows.Scan(&hash); err != nil {
			return nil, err
		}
		hashes = append(hashes, hash)
	}
	return hashes, rows.Err()
}

func (r *RefreshTokenRepository) UpdateLastUsed(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET last_used_at = NOW() WHERE token_hash = $1`
	_, err := r.DB.Exec(ctx, query, tokenHash)
	return err
}

func (r *RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := r.DB.Exec(ctx, query, tokenHash)
	return err
}

func (r *RefreshTokenRepository) DeleteAllUserTokens(ctx context.Context, userID int) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.DB.Exec(ctx, query, userID)
	return err
}

func (r *RefreshTokenRepository) CleanExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.DB.Exec(ctx, query)
	return err
}
