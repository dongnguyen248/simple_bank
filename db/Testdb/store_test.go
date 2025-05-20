package test

import (
	"context"
	"fmt"
	"testing"

	db "github.com/dongnguyen248/simple_bank/db/sqlc"
	"github.com/stretchr/testify/require"
)

const (
	numTransactions = 10
	transferAmount  = int64(10)
)

func TestTransferTxDeatlock(t *testing.T) {
	store := db.NewStore(testDB)
	account1 := createAccount(t)
	account2 := createAccount(t)
	t.Log(">> before:", account1.Balance, account2.Balance)

	errs := make(chan error)
	// run n concurrent transfer transactions
	for i := 0; i < numTransactions; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		t.Log(">> fromAccountID:", fromAccountID, "toAccountID:", toAccountID, "Amount:", transferAmount)
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			ctx := context.Background()
			_, err := store.TransferTx(ctx, db.TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        transferAmount,
			})
			errs <- err
		}()
	}

	// check results
	for i := 0; i < numTransactions; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	// check the final updated balance
	updatedAccount1, err := testQueries.GetAcount(context.Background(), account1.ID)
	require.NoError(t, err)
	t.Log(">> updatedAccount1:", updatedAccount1)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	updatedAccount2, err := testQueries.GetAcount(context.Background(), account2.ID)
	t.Log(">> updatedAccount2:", updatedAccount2)
	require.NoError(t, err)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
	t.Log(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

}

var txName = ""

func TestTransferTx(t *testing.T) {
	store := db.NewStore(testDB)
	account1 := createAccount(t)
	account2 := createAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	errs := make(chan error)
	results := make(chan db.TransferTxResult)
	// run n concurrent transfer transactions
	for i := 0; i < numTransactions; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, db.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        transferAmount,
			})
			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)
	for i := 0; i < numTransactions; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)
		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAcountID)
		require.Equal(t, transferAmount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		// check from entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -transferAmount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		// check to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, transferAmount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		// check account balance
		fmt.Println(txName, fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%transferAmount == 0) // amount must be multiple of n

		k := int(diff1 / transferAmount)
		require.True(t, k >= 1 && k <= numTransactions) // k must be between 1 and n
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	// check the final updated balance
	updatedAccount1, err := testQueries.GetAcount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAcount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(numTransactions)*transferAmount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(numTransactions)*transferAmount, updatedAccount2.Balance)

}
