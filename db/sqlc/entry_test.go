package db

import (
	"context"
	"testing"

	"github.com/begenov/backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, accound Account) Entry {
	arg := CreateEntryParams{
		AccountID: accound.ID,
		Amount:    util.RandomMany(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createTestAccaunt(t)
	createRandomEntry(t, account)
}
