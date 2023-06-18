package repository

import (
	"context"
	"database/sql"

	"github.com/begenov/backend/internal/domain"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type Account interface {
	CreateAccount(ctx context.Context, arg domain.CreateAccountParams) (domain.Account, error)
	DeleteAccount(ctx context.Context, id int) error
	GetAccount(ctx context.Context, id int) (domain.Account, error)
	ListAccounts(ctx context.Context, arg domain.ListAccountsParams) ([]domain.Account, error)
	UpdateAccount(ctx context.Context, arg domain.UpdateAccountParams) (domain.Account, error)
	GetAccountForUpdate(ctx context.Context, id int) (domain.Account, error)
	AddAccountBalance(ctx context.Context, arg domain.AddAccountBalanceParams) (domain.Account, error)
}

type Entry interface {
	CreateEntry(ctx context.Context, arg domain.CreateEntryParams) (domain.Entry, error)
	GetEntry(ctx context.Context, id int) (domain.Entry, error)
	ListEntries(ctx context.Context, arg domain.ListEntriesParams) ([]domain.Entry, error)
}

type Transfer interface {
	CreateTransfer(ctx context.Context, arg domain.CreateTransferParams) (domain.Transfer, error)
	GetTransfer(ctx context.Context, id int) (domain.Transfer, error)
	ListTransfers(ctx context.Context, arg domain.ListTransfersParams) ([]domain.Transfer, error)
}

type User interface {
	CreateUser(ctx context.Context, arg domain.CreateUserParams) (domain.User, error)
	GetUser(ctx context.Context, username string) (domain.User, error)
}

type Tx interface {
	TransferTx(ctx context.Context, arg domain.TransferTxParams) (domain.TransferTxResult, error)
}

type Repository struct {
	db       *sql.DB
	Account  Account
	Entry    Entry
	Transfer Transfer
	User     User
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db:       db,
		Account:  New(db),
		Entry:    NewEntryRepo(db),
		Transfer: NewTransferRepo(db),
		User:     NewUserRepo(db),
	}
}
