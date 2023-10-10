package db

import (
	"context"
	"testing"
	"time"

	"github.com/beatzoid/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createAndTestRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createAndTestRandomUser(t)
}

func TestGetUser(t *testing.T) {
	randomUser := createAndTestRandomUser(t)

	foundUser, err := testQueries.GetUser(context.Background(), randomUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, foundUser)

	require.Equal(t, randomUser.Username, foundUser.Username)
	require.Equal(t, randomUser.HashedPassword, foundUser.HashedPassword)
	require.Equal(t, randomUser.FullName, foundUser.FullName)
	require.Equal(t, randomUser.Email, foundUser.Email)

	require.WithinDuration(t, randomUser.PasswordChangedAt, foundUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, randomUser.CreatedAt, foundUser.CreatedAt, time.Second)
}
