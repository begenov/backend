package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/begenov/backend/internal/domain"
	"github.com/lib/pq"
)

type EntryRepo struct {
	db *sql.DB
}

func NewEntryRepo(db *sql.DB) *EntryRepo {
	return &EntryRepo{
		db: db,
	}
}

func (r *EntryRepo) CreateEntry(ctx context.Context, arg domain.CreateEntryParams) (domain.Entry, error) {
	stmt := `
		INSERT INTO entries (
			account_id,
			amount
		) VALUES (
			$1, $2
		) RETURNING id, account_id, amount, created_at
	`
	row := r.db.QueryRowContext(ctx, stmt, arg.AccountID, arg.Amount)
	var entry domain.Entry
	err := row.Scan(
		&entry.ID,
		&entry.AccountID,
		&entry.Amount,
		&entry.CreatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23503" {
			// Handle foreign key constraint violation error
			return domain.Entry{}, fmt.Errorf("invalid account ID: %w", err)
		}
		return domain.Entry{}, err
	}
	return entry, nil
}

func (r *EntryRepo) GetEntry(ctx context.Context, id int) (domain.Entry, error) {
	stmt := `SELECT "id", "account_id", "amount", "created_at" FROM "entries" WHERE "id" = $1 LIMIT 1`
	row := r.db.QueryRowContext(ctx, stmt, id)
	var i domain.Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

func (r *EntryRepo) ListEntries(ctx context.Context, arg domain.ListEntriesParams) ([]domain.Entry, error) {
	stmt := `
		SELECT id, account_id, amount, created_at FROM entries
		WHERE account_id = $1
		ORDER BY id
		LIMIT $2
		OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, stmt, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Entry
	for rows.Next() {
		var i domain.Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
