package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

// provide all functions to run db queries and transactions
type Store struct {
	// composition to extends struct functionality instead of inheritance.
	// by embedding queries inside store, all functions provided by Queries
	// will be available inside the Store struct
	*Queries
	db *sql.DB
}

// Create new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// func to generate a transaction
// function will take context and a callback function as inputs and starts a tx
// it'll create a new Queries object with that transaction and call cb and finally commit or rollback
// based on error
// functions that start with a lower letter is not exported means no other packages
// can access the execTx function
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferFundsParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferFundsResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// Performs a money transfer from one account to another
// Creates a new transfer record, add account entries and updates accounts balances
func (store *Store) TransferFundsTx(ctx context.Context, arg TransferFundsParams) (TransferFundsResult, error) {
	var result TransferFundsResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)
		fmt.Println(txName, "Create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "Create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // money is moving out of this account
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "Create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// fmt.Println(txName, "Get account 1 for update")
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)

		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName, "Update account 1 balance")
		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.FromAccountID,
		// 	Balance: account1.Balance - arg.Amount,
		// })
		result.FromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}

		// fmt.Println(txName, "Get account 1 for update")
		// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)

		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName, "Update account 2 balance")
		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.ToAccountID,
		// 	Balance: account2.Balance + arg.Amount,
		// })

		result.ToAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
