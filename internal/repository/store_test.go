package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/begenov/backend/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewRepository(db)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	n := 5
	amount := 10

	errs := make(chan error)
	results := make(chan domain.TransferTxResult)
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, domain.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)
		results := <-results
		require.NotEmpty(t, results)

		transfer := results.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.ID)

		_, err = store.Transfer.GetTransfer(ctx, transfer.ID)
		require.NoError(t, err)

		fromEntry := results.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.Entry.GetEntry(ctx, fromEntry.ID)
		require.NoError(t, err)

		toEntry := results.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)
		require.NotZero(t, toEntry.ID)

		_, err = store.Entry.GetEntry(ctx, toEntry.ID)
		require.NoError(t, err)

		fromAccount := results.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := results.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// fmt.Println("check balance:", toAccount.Balance, fromAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := diff1 / amount
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updateAccount1, err := store.Account.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := store.Account.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">>after:", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance-n*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+n*amount, updateAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewRepository(db)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	n := 40
	amount := 10
	errs := make(chan error)

	for i := 0; i < n; i++ {

		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {

			_, err := store.TransferTx(ctx, domain.TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updateAccount1, err := store.Account.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := store.Account.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">>after:", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
