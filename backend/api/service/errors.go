package service

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
)

// mapNotFound converts a driver-level "no rows" error into the given domain
// sentinel; any other error is returned unchanged so it surfaces as a 500.
func mapNotFound(err error, sentinel error) error {
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		return sentinel
	}
	return err
}

// ValidationError carries a user-facing validation message across layers so
// handlers can surface it as a 400 without string matching.
type ValidationError struct{ Msg string }

func (e *ValidationError) Error() string { return e.Msg }
