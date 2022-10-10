package sqlc

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferFundsTx(t *testing.T) {
	store := NewStore(testDB)
	createAccount1Args := getRandomTestAccountParams()
	createAccount2Args := getRandomTestAccountParams()
	account1, account1Err := createTestAccount(createAccount1Args)
	account2, account2Err := createTestAccount(createAccount2Args)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	if account1Err != nil {
		log.Fatal("Failed to create test account1")
	}

	if account2Err != nil {
		log.Fatal("Failed to create test account2")
	}

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferFundsResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		// use go keyword to start a go routine
		/*
			* we cannot use testify directly here because this is running on a separate go rountine
			than TransferFundsTx is on. Since this function is running inside a different routine,
			there is no guarantee it will stop if the whole test if a condition is not sastified.
			The correct way to verify the result and error is to send the output back to the main
			go routine. To do that we use channels.
			Channels is designed to connect concurrent go routines and allow them to safely share
			data with each other without explict locking. In this case, we need one channel to receive
			errors and one channel to receive transaction results
		*/
		go func() {
			result, err := store.TransferFundsTx(context.Background(), TransferFundsParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		// receiving all the errors from the errs channel
		// the variable on the left, stores the receives data and the arrow is on the left of the
		// channel to send the data
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransferById(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntryById(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntryById(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balances
		fmt.Println(">> transaction:", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		/**
		the amount should be divisible by the amount in each tx and be positive
		the reason is the balance will be decreased with each transaction
		*/
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := store.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
