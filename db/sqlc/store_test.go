package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	fmt.Printf(">> before : %d, %d \n", account1.Balance, account2.Balance)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx%d\n", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	//check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check accounts
		fromAccount := result.FromAccount
		_, err = store.GetAccount(context.Background(), fromAccount.ID)
		require.NoError(t, err)
		require.NotEmpty(t, fromAccount)

		toAccount := result.ToAccount
		_, err = store.GetAccount(context.Background(), toAccount.ID)
		require.NoError(t, err)
		require.NotEmpty(t, toAccount)

		//check accounts' balance
		fmt.Printf(">> tx : %d, %d \n", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		k := int(diff1 / amount)
		require.Equal(t, diff1, diff2)
		require.True(t, diff1%amount == 0)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	updatedAcocunt1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAcocunt2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAcocunt1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAcocunt2.Balance)
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("account1 balance : ", account1.Balance, "account2 balance : ", account2.Balance)
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromID := account1.ID
		toID := account2.ID

		if i%2 == 1 {
			fromID = account2.ID
			toID = account1.ID
		}
		fmt.Printf(" init i: %d \n", i)
		go func() {
			txName := fmt.Sprintf("tx%d\n", i+1)
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromID,
				ToAccountID:   toID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAcocunt1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAcocunt2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAcocunt1.Balance)
	require.Equal(t, account2.Balance, updatedAcocunt2.Balance)
}
