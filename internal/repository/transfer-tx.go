package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/begenov/backend/internal/domain"
)

func (r *Repository) TransferTx(ctx context.Context, arg domain.TransferTxParams) (domain.TransferTxResult, error) {
	var result domain.TransferTxResult

	err := r.execTx(ctx, func(account Account, transfer Transfer, entry Entry) error {
		var err error

		result.Transfer, err = transfer.CreateTransfer(ctx, domain.CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = entry.CreateEntry(ctx, domain.CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = entry.CreateEntry(ctx, domain.CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
		// TODO: update accounts' balance
	})

	return result, err
}

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
