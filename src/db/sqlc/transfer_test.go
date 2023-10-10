package db

import (
	"context"
	"testing"

	"github.com/beatzoid/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createAndTestRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	randomAccount1 := createAndTestRandomAccount(t)
	randomAccount2 := createAndTestRandomAccount(t)

	createAndTestRandomTransfer(t, randomAccount1, randomAccount2)
}

func TestGetTransfer(t *testing.T) {
	randomAccount1 := createAndTestRandomAccount(t)
	randomAccount2 := createAndTestRandomAccount(t)
	randomTransfer := createAndTestRandomTransfer(t, randomAccount1, randomAccount2)

	foundTransfer, err := testQueries.GetTransfer(context.Background(), randomTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, foundTransfer)

	require.Equal(t, randomTransfer.ID, foundTransfer.ID)
	require.Equal(t, randomTransfer.FromAccountID, foundTransfer.FromAccountID)
	require.Equal(t, randomTransfer.ToAccountID, foundTransfer.ToAccountID)
	require.Equal(t, randomTransfer.Amount, foundTransfer.Amount)
	require.WithinDuration(t, randomTransfer.CreatedAt, foundTransfer.CreatedAt, 0)
}

func TestListTransfer(t *testing.T) {
	randomAccount1 := createAndTestRandomAccount(t)
	randomAccount2 := createAndTestRandomAccount(t)

	for i := 0; i < 5; i++ {
		createAndTestRandomTransfer(t, randomAccount1, randomAccount2)
		createAndTestRandomTransfer(t, randomAccount2, randomAccount1)
	}

	arg := ListTransfersParams{
		FromAccountID: randomAccount1.ID,
		ToAccountID:   randomAccount1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == randomAccount1.ID || transfer.ToAccountID == randomAccount1.ID)
	}

	arg = ListTransfersParams{
		FromAccountID: randomAccount1.ID,
		ToAccountID:   randomAccount1.ID,
		Limit:         -5,
		Offset:        -5,
	}

	_, err = testQueries.ListTransfers(context.Background(), arg)

	require.Error(t, err)
}
