package repository

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/pkg/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func init() {
	NewRepoTest()
}

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var ctx context.Context = context.Background()

var repo *AccountRepo

var db *sql.DB

func NewRepoTest() {
	var err error
	db, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	transferRepo = NewTransferRepo(db)
	repo = New(db)
	entryRepo = NewEntryRepo(db)
}

func TestCreateACcount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	account2, err := repo.GetAccount(ctx, account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := domain.UpdateAccountParams{
		ID:      account1.ID,
		Balance: int(util.RandomMany()),
	}
	account2, err := repo.UpdateAccount(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestDeleateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := repo.DeleteAccount(ctx, account1.ID)
	require.NoError(t, err)

	account2, err := repo.GetAccount(ctx, account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := domain.ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := repo.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, arg.Limit)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}

func createRandomAccount(t *testing.T) domain.Account {
	arg := domain.CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  int(util.RandomMany()),
		Currency: util.RandomCurrency(),
	}

	account, err := repo.CreateAccount(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
