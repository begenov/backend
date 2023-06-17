package service

import (
	"context"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/internal/repository"
)

type AccountService struct {
	repo repository.Account
}

func NewAccountService(repo repository.Account) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, arg domain.CreateAccountParams) (domain.Account, error) {
	return s.repo.CreateAccount(ctx, arg)
}

func (s *AccountService) GetAccountByID(ctx context.Context, id int) (domain.Account, error) {
	return s.repo.GetAccount(ctx, id)
}

func (s *AccountService) ListAccounts(ctx context.Context, arg domain.ListAccountsParams) ([]domain.Account, error) {
	return s.repo.ListAccounts(ctx, arg)
}
