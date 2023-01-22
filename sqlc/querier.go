// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package sqlc

import (
	"context"
)

type Querier interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccountById(ctx context.Context, id int64) (Account, error)
	// tells postgres that we dont update the key of id column of the account table to avoid deadlock
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	GetAccounts(ctx context.Context) ([]Account, error)
	GetEntries(ctx context.Context, arg GetEntriesParams) ([]Entry, error)
	GetEntryById(ctx context.Context, id int64) (Entry, error)
	GetTransferById(ctx context.Context, id int64) (Transfer, error)
	GetTransfers(ctx context.Context, arg GetTransfersParams) ([]Transfer, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
	UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) (Account, error)
}

var _ Querier = (*Queries)(nil)
