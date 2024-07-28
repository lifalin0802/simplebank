package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account1)

	require.Equal(t, account2.ID, account1.ID)
	require.Equal(t, account2.Currency, account1.Currency)
	require.Equal(t, account2.Owner, account1.Owner)
	require.Equal(t, account2.Balance, account1.Balance)

	require.WithinDuration(t, account1.CaretedAt, account2.CaretedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account2.ID, account1.ID)
	require.Equal(t, account2.Currency, account1.Currency)
	require.Equal(t, account2.Owner, account1.Owner)
	require.NotEqual(t, account2.Balance, account1.Balance)
	require.WithinDuration(t, account1.CaretedAt, account2.CaretedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err2 := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err2)
	require.EqualError(t, err2, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}

func TestListAccount(t *testing.T) {
	createRandomAccount(t)
	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		Limit:  1,
		Offset: 0,
	})

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
