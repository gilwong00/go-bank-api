package sqlc

import (
	"context"
	"go-bank-api/pkg/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsertIntoAcount(t *testing.T) {
	arg := InsertIntoAccountParams{
		Owner:    util.GetRandomOwner(),
		Balance:  util.GetRandomBalance(),
		Currency: util.GetCurrencyType(),
	}

	account, err := testQueries.InsertIntoAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
