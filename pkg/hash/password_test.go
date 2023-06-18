package hash

import (
	"testing"

	"github.com/begenov/backend/pkg/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var hash PasswordHasher

func init() {
	hash = NewHash()
}

func TestPassword(t *testing.T) {
	password := util.RandomString(6)

	hashedPassword, err := hash.GenerateFromPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = hash.CompareHashAndPassword(hashedPassword, password)
	require.NoError(t, err)

	wrongPassword := util.RandomString(6)
	err = hash.CompareHashAndPassword(hashedPassword, wrongPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := hash.GenerateFromPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword, hashedPassword2)
}

func TestWrongPassword(t *testing.T) {
	password := util.RandomString(6)

	hashedPassword, err := hash.GenerateFromPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	wrongPassword := util.RandomString(6)
	err = hash.CompareHashAndPassword(hashedPassword, wrongPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
