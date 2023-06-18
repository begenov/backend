package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/begenov/backend/internal/domain"
)

var txKey = struct{}{}

func (r *Repository) TransferTx(ctx context.Context, arg domain.TransferTxParams) (domain.TransferTxResult, error) {
	var result domain.TransferTxResult

	err := r.execTx(ctx, func(account Account, transfer Transfer, entry Entry) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		result.Transfer, err = transfer.CreateTransfer(ctx, domain.CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = entry.CreateEntry(ctx, domain.CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")

		result.ToEntry, err = entry.CreateEntry(ctx, domain.CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = r.addMoney(ctx, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = r.addMoney(ctx, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

		}

		fmt.Println(txName, "updated account 2 balance:", result.ToAccount.Balance)
		log.Println(result.FromAccount.ID == arg.FromAccountID, result.ToAccount.ID == arg.ToAccountID)
		return err
	})

	return result, err
}

func (r *Repository) addMoney(ctx context.Context, fromAccountID int, fromAmount int, toAccountID int, toAmount int) (account1 domain.Account, account2 domain.Account, err error) {
	account1, err = r.Account.AddAccountBalance(ctx, domain.AddAccountBalanceParams{
		ID:     fromAccountID,
		Amount: fromAmount,
	})

	if err != nil {
		return
	}
	account2, err = r.Account.AddAccountBalance(ctx, domain.AddAccountBalanceParams{
		ID:     toAccountID,
		Amount: toAmount,
	})

	return
}
