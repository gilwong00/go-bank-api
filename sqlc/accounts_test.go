package sqlc

import (
	"context"
	"go-bank-api/pkg/util"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createTestAccount(arg CreateAccountParams) (Account, error) {
	return testQueries.CreateAccount(context.Background(), arg)
}

func getRandomTestAccountParams() CreateAccountParams {
	arg := CreateAccountParams{
		Owner:    util.GetRandomOwner(),
		Balance:  util.GetRandomBalance(),
		Currency: util.GetCurrencyType(),
	}

	return arg
}

func TestCreateAccount(t *testing.T) {
	arg := getRandomTestAccountParams()
	account, err := createTestAccount(arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccountById(t *testing.T) {
	arg := getRandomTestAccountParams()
	account1, createErr := createTestAccount(arg)

	if createErr != nil {
		log.Fatal("Cannot create test account")
	}

	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
}

func TestUpdateAccount(t *testing.T) {
	createArg := getRandomTestAccountParams()
	account1, err := createTestAccount(createArg)

	if err != nil {
		log.Fatal("Error creating account")
	}

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.GetRandomBalance(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
