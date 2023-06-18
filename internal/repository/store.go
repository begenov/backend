package repository

import (
	"context"
	"database/sql"
	"fmt"
)

func (r *Repository) execTx(ctx context.Context, fn func(Account, Transfer, Entry) error) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = fn(r.Account, r.Transfer, r.Entry)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
