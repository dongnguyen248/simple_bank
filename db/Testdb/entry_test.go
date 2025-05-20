package test

import (
	"context"
	"testing"

	db "github.com/dongnguyen248/simple_bank/db/sqlc"
	"github.com/dongnguyen248/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func CreateEntry(t *testing.T) db.Entry {
	account := createAccount(t)
	// Create an entry for the account
	// with a random amount
	arg := db.CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

func TestCreateEntry(t *testing.T) {
	entry := CreateEntry(t)
	require.NotEmpty(t, entry)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	require.Equal(t, entry.AccountID, entry.AccountID)
	require.Equal(t, entry.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
}

func TestGetEntry(t *testing.T) {
	entry1 := CreateEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, 0)
}
func TestUpdateEntry(t *testing.T) {
	entry := CreateEntry(t)
	arg := db.UpdateEntryParams{
		ID:     entry.ID,
		Amount: util.RandomMoney(),
	}
	entry2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry.ID, entry2.ID)
	require.Equal(t, entry.AccountID, entry2.AccountID)
	require.Equal(t, arg.Amount, entry2.Amount)
}
func TestDeleteEntry(t *testing.T) {
	entry := CreateEntry(t)
	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	// Check that the account has been deleted
	_, err = testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
}

func TestListEntries(t *testing.T) {
	// Create 10 random entries
	for i := 0; i < 10; i++ {
		CreateEntry(t)
	}

	arg := db.ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
