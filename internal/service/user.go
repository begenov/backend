package service

import (
	"context"
	"time"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/internal/repository"
	"github.com/begenov/backend/pkg/auth"
	"github.com/begenov/backend/pkg/e"
	"github.com/begenov/backend/pkg/hash"
)

type UserService struct {
	repo                repository.User
	hash                hash.PasswordHasher
	token               auth.TokenManager
	accessTokenDuration time.Duration
}

func NewUserService(repo repository.User, hash hash.PasswordHasher, token auth.TokenManager, accessTokenDuration time.Duration) *UserService {
	return &UserService{
		repo:                repo,
		hash:                hash,
		token:               token,
		accessTokenDuration: accessTokenDuration,
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

func (s *UserService) GetUserByUsername(ctx context.Context, username string, password string) (domain.LoginUserResponse, error) {

	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		return domain.LoginUserResponse{}, err
	}
	err = s.hash.CompareHashAndPassword(user.HashedPassword, password)
	if err != nil {
		return domain.LoginUserResponse{}, e.ErrPassword
	}
	accessToken, err := s.token.NewJWT(user.Username, s.accessTokenDuration)
	if err != nil {
		return domain.LoginUserResponse{}, e.ErrInvalidToken
	}

	response := domain.LoginUserResponse{
		AccessToken: accessToken,
		User: domain.UserResponse{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: user.PasswordChangedAt,
			CreatedAt:         user.CreatedAt,
		},
	}
	return response, nil
}
