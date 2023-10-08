package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// TransferTxParams contains the input parameters for the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, transactionFunction func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	// Create a new query object to interact with the DB
	queries := New(tx)
	// Execute the passed in transaction
	err = transactionFunction(queries)

	if err != nil {
		// Rollback the transaction (undo the changes) if it errors
		if rollbackError := tx.Rollback(); rollbackError != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rollbackError)
		}

		return err
	}

	// Commit the transaction (save the changes) if there are no errors
	return tx.Commit()
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add the accounts entries,
// and updates the account balances with a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create a transfer record showing the money leaving one account and entering the other
		// CreateTransferParams(arg) is the same as
		// CreateTransferParams{
		//	FromAccountID: arg.FromAccountID,
		//	ToAccountID:   arg.ToAccountID,
		//	Amount:        arg.Amount,
		// }
		// But it makes the linter happy (and looks a little cleaner)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))

		if err != nil {
			return err
		}

		// Create an entry in the FromAccount showing the amount leaving
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		// Create an entry in the ToAccount showing the amount entering
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		//  Update account balances

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	if err != nil {
		return
	}

	return
}
