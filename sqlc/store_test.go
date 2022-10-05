package sqlc

import (
	"context"
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

	if account1Err != nil {
		log.Fatal("Failed to create test account1")
	}

	if account2Err != nil {
		log.Fatal("Failed to create test account2")
	}

	// run n amount of current transactions
	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferFundsResult)

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

			// send an error to the errors channel
			errs <- err

			// send results to the results channel
			results <- result
		}()
	}

	// check results
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
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// check to see if record is property created in the DB
		_, err = store.GetTransferById(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check from entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntryById(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check to entries
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntryById(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO check balances
	}
}
