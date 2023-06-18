package service

import (
	"context"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/internal/repository"
	"github.com/begenov/backend/pkg/hash"
)

type UserService struct {
	repo repository.User
	hash hash.PasswordHasher
}

func NewUserService(repo repository.User, hash hash.PasswordHasher) *UserService {
	return &UserService{
		repo: repo,
		hash: hash,
	}
}

func (s *UserService) CreateUser(ctx context.Context, arg domain.CreateUserParams) (domain.User, error) {
	var err error
	arg.HashedPassword, err = s.hash.GenerateFromPassword(arg.HashedPassword)
	if err != nil {
		return domain.User{}, err
	}
	return s.repo.CreateUser(ctx, arg)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	return s.repo.GetUser(ctx, username)
}
