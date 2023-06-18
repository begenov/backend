package repository

import (
	"testing"
	"time"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/pkg/util"
	"github.com/stretchr/testify/require"
)

var userRepo *UserRepo

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := userRepo.GetUser(ctx, user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func createRandomUser(t *testing.T) domain.User {
	// hashedPassword, err := h.GenerateFromPassword(util.RandomString(6))
	// require.NoError(t, err)

	arg := domain.CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "hashedPassword",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := userRepo.CreateUser(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())
	return user
}
