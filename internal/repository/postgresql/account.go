package postgresql

import (
	"context"
	"database/sql"

	"github.com/begenov/backend/internal/domain"
)

type AccountRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *AccountRepo {
	return &AccountRepo{
		db: db,
	}
}

func (r *AccountRepo) CreateAccount(ctx context.Context, arg domain.CreateAccountParams) (domain.Account, error) {
	stmt := `INSERT INTO accounts (
		owner, 
		balance, 
		currency
	) VALUES (
		$1, $2, $3
	) RETURNING id, owner, balance, currency, created_at
	`

	row := r.db.QueryRowContext(ctx, stmt, arg.Owner, arg.Balance, arg.Currency)
	var i domain.Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)

	return i, err
}

func (r *AccountRepo) DeleteAccount(ctx context.Context, id int) error {
	stmt := `DELETE FROM accounts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, stmt, id)
	return err
}

func (r *AccountRepo) GetAccount(ctx context.Context, id int) (domain.Account, error) {
	stmt := `SELECT id, owner, balance, currency, created_at FROM accounts
	WHERE id = $1 LIMIT 1`
	row := r.db.QueryRowContext(ctx, stmt, id)
	var i domain.Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

func (r *AccountRepo) ListAccounts(ctx context.Context, arg domain.ListAccountsParams) ([]domain.Account, error) {
	stmt := `SELECT id, owner, balance, currency, created_at FROM accounts
	ORDER BY id
	LIMIT $1
	OFFSET $2`
	row, err := r.db.QueryContext(ctx, stmt, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	var items []domain.Account
	for row.Next() {
		var i domain.Account
		if err := row.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err := row.Close(); err != nil {
		return nil, err
	}
	if err := row.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *AccountRepo) UpdateAccount(ctx context.Context, arg domain.UpdateAccountParams) (domain.Account, error) {
	stmt := `UPDATE accounts
	SET balance = $2
	WHERE id = $1
	RETURNING id, owner, balance, currency, created_at`

	row := r.db.QueryRowContext(ctx, stmt, arg.ID, arg.Balance)
	var i domain.Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}
