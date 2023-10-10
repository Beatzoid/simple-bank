package db

import (
	"context"
	"testing"
	"time"

	"github.com/beatzoid/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createAndTestRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	randomAccount := createAndTestRandomAccount(t)
	createAndTestRandomEntry(t, randomAccount)
}

func TestGetEntry(t *testing.T) {
	randomAccount := createAndTestRandomAccount(t)
	randomEntry := createAndTestRandomEntry(t, randomAccount)

	foundEntry, err := testQueries.GetEntry(context.Background(), randomEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, foundEntry)

	require.Equal(t, randomEntry.ID, foundEntry.ID)
	require.Equal(t, randomEntry.AccountID, foundEntry.AccountID)
	require.Equal(t, randomEntry.Amount, foundEntry.Amount)
	require.WithinDuration(t, randomEntry.CreatedAt, foundEntry.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	randomAccount := createAndTestRandomAccount(t)

	for i := 0; i < 10; i++ {
		createAndTestRandomEntry(t, randomAccount)
	}

	arg := ListEntriesParams{
		AccountID: randomAccount.ID,
		Limit:     5,
		Offset:    5, // Skip the first 5 records and return the next 5
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}

	arg = ListEntriesParams{
		AccountID: randomAccount.ID,
		Limit:     -5,
		Offset:    -5, // Skip the first 5 records and return the next 5
	}

	_, err = testQueries.ListEntries(context.Background(), arg)

	require.Error(t, err)
}
