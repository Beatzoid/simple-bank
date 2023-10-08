package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	randomAccount1 := createAndTestRandomAccount(t)
	randomAccount2 := createAndTestRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// Run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: randomAccount1.ID,
				ToAccountID:   randomAccount2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer

		// Check transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, randomAccount1.ID, transfer.FromAccountID)
		require.Equal(t, randomAccount2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries

		fromEntry := result.FromEntry

		require.NotEmpty(t, fromEntry)
		require.Equal(t, randomAccount1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)
		require.Equal(t, randomAccount2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check accounts

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, randomAccount1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, randomAccount2.ID, toAccount.ID)

		// Check accounts balance

		// The amount of money going out of account 1
		diff1 := randomAccount1.Balance - fromAccount.Balance
		// The amount of money going into account 2
		diff2 := toAccount.Balance - randomAccount2.Balance

		require.Equal(t, diff1, diff2)
		// Cannot transfer a negative amount of money
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount... n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the final updated balance of the two accounts

	updatedAccount1, err := testQueries.GetAccount(context.Background(), randomAccount1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), randomAccount2.ID)
	require.NoError(t, err)

	// Check the final updated balance of the two accounts

	// Verify the amount of money that should be going out of account1 actually went out
	require.Equal(t, randomAccount1.Balance-int64(n)*amount, updatedAccount1.Balance)

	// Verify the amount of money that should be going out into account2 actually went in
	require.Equal(t, randomAccount2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	randomAccount1 := createAndTestRandomAccount(t)
	randomAccount2 := createAndTestRandomAccount(t)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	// Run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		fromAccountID := randomAccount1.ID
		toAccountID := randomAccount2.ID

		if i%2 == 1 {
			fromAccountID = randomAccount2.ID
			toAccountID = randomAccount1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Check the final updated balance of the two accounts

	updatedAccount1, err := testQueries.GetAccount(context.Background(), randomAccount1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), randomAccount2.ID)
	require.NoError(t, err)

	// The amount of money in both accounts should stay the same
	// because we are just transferring the money back and forth
	// between the two accounts
	require.Equal(t, randomAccount1.Balance, updatedAccount1.Balance)
	require.Equal(t, randomAccount2.Balance, updatedAccount2.Balance)
}
