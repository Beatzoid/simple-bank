package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/beatzoid/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createAndTestRandomAccount(t *testing.T) Account {
	user := createAndTestRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createAndTestRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	randomAccount := createAndTestRandomAccount(t)

	foundAccount, err := testQueries.GetAccount(context.Background(), randomAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, foundAccount)

	require.Equal(t, randomAccount.ID, foundAccount.ID)
	require.Equal(t, randomAccount.Owner, foundAccount.Owner)
	require.Equal(t, randomAccount.Balance, foundAccount.Balance)
	require.Equal(t, randomAccount.Currency, foundAccount.Currency)
	require.WithinDuration(t, randomAccount.CreatedAt, foundAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	randomAccount := createAndTestRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      randomAccount.ID,
		Balance: util.RandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, randomAccount.ID, updatedAccount.ID)
	require.Equal(t, randomAccount.Owner, updatedAccount.Owner)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
	require.Equal(t, randomAccount.Currency, updatedAccount.Currency)
	require.WithinDuration(t, randomAccount.CreatedAt, updatedAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	randomAccount := createAndTestRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), randomAccount.ID)

	require.NoError(t, err)

	foundAccount, err := testQueries.GetAccount(context.Background(), randomAccount.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, foundAccount)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account

	for i := 0; i < 10; i++ {
		lastAccount = createAndTestRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
