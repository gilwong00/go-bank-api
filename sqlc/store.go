package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

// mock store interface for testing
type Store interface {
	Querier // generated from sqlc: emit_interface:true
	TransferFundsTx(ctx context.Context, arg TransferFundsParams) (TransferFundsResult, error)
}

// provide all functions to run db queries and transactions
type DbStore struct {
	// composition to extends struct functionality instead of inheritance.
	// by embedding queries inside store, all functions provided by Queries
	// will be available inside the Store struct
	*Queries
	db *sql.DB
}

// Create new store
func NewStore(db *sql.DB) Store {
	return &DbStore{
		db:      db,
		Queries: New(db),
	}
}

// func to generate a transaction
// function will take context and a callback function as inputs and starts a tx
// it'll create a new Queries object with that transaction and call cb and finally commit or rollback
// based on errorfunctions that start with a lower letter is not exported means no other packages
// can access the execTx function
func (store *DbStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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

// Performs a money transfer from one account to another
// Creates a new transfer record, add account entries and updates accounts balances
func (store *DbStore) TransferFundsTx(ctx context.Context, arg TransferFundsParams) (TransferFundsResult, error) {
	var result TransferFundsResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // money is moving out of this account
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		/*
			Run updates on sequential Ids first to avoid deadlocks
		*/
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoneyToBalance(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoneyToBalance(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return err
	})

	return result, err
}

func addMoneyToBalance(
	ctx context.Context,
	q *Queries,
	account1Id int64,
	amount1 int64,
	account2Id int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     account1Id,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     account2Id,
		Amount: amount2,
	})

	return
}
