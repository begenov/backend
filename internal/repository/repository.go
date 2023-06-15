package repository

import (
	"context"
	"database/sql"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/internal/repository/postgresql"
)

type Account interface {
	CreateAccount(ctx context.Context, arg domain.CreateAccountParams) (domain.Account, error)
	DeleteAccount(ctx context.Context, id int) error
	GetAccount(ctx context.Context, id int) (domain.Account, error)
	ListAccounts(ctx context.Context, arg domain.ListAccountsParams) ([]domain.Account, error)
	UpdateAccount(ctx context.Context, arg domain.UpdateAccountParams) (domain.Account, error)
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

type Repository struct {
	Account  Account
	Entry    Entry
	Transfer Transfer
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Account:  postgresql.New(db),
		Entry:    postgresql.NewEntryRepo(db),
		Transfer: postgresql.NewTransferRepo(db),
	}
}
