package database

import (
	"context"
	"database/sql"
	"fmt"
)

// UnitOfWork owns transaction boundaries so repositories never commit on their own.
type UnitOfWork struct {
	db *sql.DB
}

func NewUnitOfWork(db *sql.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

// Do runs fn inside a transaction and commits only when fn returns nil.
// Do not nest calls and do not share the executor across goroutines.
func (u *UnitOfWork) Do(ctx context.Context, fn func(Executor) error) error {
	return u.DoWithOptions(ctx, nil, fn)
}

func (u *UnitOfWork) DoWithOptions(ctx context.Context, opts *sql.TxOptions, fn func(Executor) error) error {
	tx, err := u.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("database: begin transaction: %w", err)
	}
	// Rollback after a successful commit is a no-op, so this covers every exit path.
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("database: commit transaction: %w", err)
	}

	return nil
}
