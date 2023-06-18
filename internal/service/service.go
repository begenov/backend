package service

import (
	"context"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/internal/repository"
	"github.com/begenov/backend/pkg/hash"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Account interface {
	CreateAccount(ctx context.Context, arg domain.CreateAccountParams) (domain.Account, error)
	GetAccountByID(ctx context.Context, id int) (domain.Account, error)
	ListAccounts(ctx context.Context, arg domain.ListAccountsParams) ([]domain.Account, error)
}

type TransferTx interface {
	TransferTx(ctx context.Context, arg domain.TransferTxParams) (domain.TransferTxResult, error)
}

type User interface {
	CreateUser(ctx context.Context, arg domain.CreateUserParams) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}

type Service struct {
	Account    Account
	TransferTx TransferTx
	User       User
}

func NewService(repo *repository.Repository, hash hash.PasswordHasher) *Service {
	return &Service{
		Account:    NewAccountService(repo.Account),
		TransferTx: NewTransferService(repo),
		User:       NewUserService(repo.User, hash),
	}
}
