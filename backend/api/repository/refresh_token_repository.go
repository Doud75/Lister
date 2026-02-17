package repository

import (
	"context"
	"setlist/api/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, tokenHash string) (userID int, expiresAt time.Time, err error)
	GetUserTokenHashes(ctx context.Context, userID int) ([]string, error)
	UpdateLastUsed(ctx context.Context, tokenHash string) error
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	DeleteAllUserTokens(ctx context.Context, userID int) error
	ReplaceUserRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error
	CleanExpiredTokens(ctx context.Context) error
	GetAllValidTokens(ctx context.Context) ([]model.RefreshToken, error)
}

type PgRefreshTokenRepository struct {
	DB *pgxpool.Pool
}

func (r *PgRefreshTokenRepository) StoreRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (r *PgRefreshTokenRepository) FindRefreshToken(ctx context.Context, tokenHash string) (userID int, expiresAt time.Time, err error) {
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash = $1`
	err = r.DB.QueryRow(ctx, query, tokenHash).Scan(&userID, &expiresAt)
	return
}

func (r *PgRefreshTokenRepository) GetUserTokenHashes(ctx context.Context, userID int) ([]string, error) {
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

func (r *PgRefreshTokenRepository) UpdateLastUsed(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET last_used_at = NOW() WHERE token_hash = $1`
	_, err := r.DB.Exec(ctx, query, tokenHash)
	return err
}

func (r *PgRefreshTokenRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := r.DB.Exec(ctx, query, tokenHash)
	return err
}

func (r *PgRefreshTokenRepository) DeleteAllUserTokens(ctx context.Context, userID int) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.DB.Exec(ctx, query, userID)
	return err
}

func (r *PgRefreshTokenRepository) ReplaceUserRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	insertQuery := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, insertQuery, userID, tokenHash, expiresAt)
	if err != nil {
		return err
	}

	deleteQuery := `
		DELETE FROM refresh_tokens 
		WHERE user_id = $1 
		AND id NOT IN (
			SELECT id FROM refresh_tokens 
			WHERE user_id = $1 
			ORDER BY created_at DESC 
			LIMIT 3
		)`
	_, err = tx.Exec(ctx, deleteQuery, userID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PgRefreshTokenRepository) CleanExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.DB.Exec(ctx, query)
	return err
}

func (r *PgRefreshTokenRepository) GetAllValidTokens(ctx context.Context) ([]model.RefreshToken, error) {
	var tokens []model.RefreshToken
	query := `SELECT user_id, token_hash, expires_at FROM refresh_tokens WHERE expires_at > NOW()`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var token model.RefreshToken
		if err := rows.Scan(&token.UserID, &token.TokenHash, &token.ExpiresAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, rows.Err()
}
