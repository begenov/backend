package repository

import (
	"context"
	"database/sql"

	"github.com/begenov/backend/internal/domain"
)

type TransferRepo struct {
	db *sql.DB
}

func NewTransferRepo(db *sql.DB) *TransferRepo {
	return &TransferRepo{
		db: db,
	}
}

func (r *TransferRepo) CreateTransfer(ctx context.Context, arg domain.CreateTransferParams) (domain.Transfer, error) {
	stmt := `INSERT INTO transfers (
		from_account_id,
		to_account_id,
		amount
	) VALUES (
		$1, $2, $3
	) RETURNING id, from_account_id, to_account_id, amount, "created_at"`
	row := r.db.QueryRowContext(ctx, stmt, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i domain.Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

func (r *TransferRepo) GetTransfer(ctx context.Context, id int) (domain.Transfer, error) {
	stmt := `SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
	WHERE id = $1 LIMIT 1`
	row := r.db.QueryRowContext(ctx, stmt, id)
	var i domain.Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

func (r *TransferRepo) ListTransfers(ctx context.Context, arg domain.ListTransfersParams) ([]domain.Transfer, error) {
	stmt := `SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
	WHERE 
		from_account_id = $1 OR
		to_account_id = $2
	ORDER BY id
	LIMIT $3
	OFFSET $4`
	rows, err := r.db.QueryContext(ctx, stmt, arg.FromAccountID, arg.ToAccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	var items []domain.Transfer
	for rows.Next() {
		var i domain.Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}
