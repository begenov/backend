package service

import (
	"context"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/internal/repository"
)

type Account interface {
	CreateAccount(ctx context.Context, arg domain.CreateAccountParams) (domain.Account, error)
	GetAccountByID(ctx context.Context, id int) (domain.Account, error)
	ListAccounts(ctx context.Context, arg domain.ListAccountsParams) ([]domain.Account, error)
}

type Service struct {
	Account Account
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Account: NewAccountService(repo.Account),
	}
}
