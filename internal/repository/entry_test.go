package repository

import (
	"testing"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/pkg/util"
	"github.com/stretchr/testify/require"
)

var entryRepo *EntryRepo

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	entry2, err := entryRepo.GetEntry(ctx, entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	require.Equal(t, entry1.ID, entry2.ID)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := domain.ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := entryRepo.ListEntries(ctx, arg)
	require.NoError(t, err)
	require.Len(t, entries, arg.Limit)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}

}

func createRandomEntry(t *testing.T, account domain.Account) domain.Entry {
	arg := domain.CreateEntryParams{
		AccountID: account.ID,
		Amount:    int(util.RandomMany()),
	}

	entry, err := entryRepo.CreateEntry(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}
