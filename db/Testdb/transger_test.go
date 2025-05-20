package test

import (
	"context"
	"testing"

	db "github.com/dongnguyen248/simple_bank/db/sqlc"
	"github.com/dongnguyen248/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func CreateTransfer(t *testing.T) db.Transfer {
	account1 := createAccount(t)
	account2 := createAccount(t)
	// Create a transfer between the two accounts
	arg := db.CreateTransferParams{
		FromAccountID: account1.ID,
		ToAcountID:    account2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAcountID, transfer.ToAcountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	transfer := CreateTransfer(t)
	require.NotEmpty(t, transfer)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	require.Equal(t, transfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, transfer.ToAcountID, transfer.ToAcountID)
	require.Equal(t, transfer.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := CreateTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAcountID, transfer2.ToAcountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, 0)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := CreateTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err)
	require.Empty(t, transfer2)
}

func TestUpdateTransfer(t *testing.T) {
	transfer := CreateTransfer(t)
	arg := db.UpdateTransferParams{
		ID:     transfer.ID,
		Amount: util.RandomMoney(),
	}
	transfer2, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer.ID, transfer2.ID)
	require.Equal(t, transfer.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, arg.Amount, transfer2.Amount)
	require.Equal(t, transfer.ToAcountID, transfer2.ToAcountID)
	require.WithinDuration(t, transfer.CreatedAt, transfer2.CreatedAt, 0)
}
