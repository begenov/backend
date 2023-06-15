package repository

import (
	"context"
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
		go func() {
			result, err := store.TransferTx(context.Background(), domain.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		results := <-results
		require.NotEmpty(t, results)

		transfer := results.Transfer
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

		// TODO: check accounts' balance
	}
}
