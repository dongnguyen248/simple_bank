package test

import (
	"context"
	"testing"

	db "github.com/dongnguyen248/simple_bank/db/sqlc"
	"github.com/dongnguyen248/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createUser(t *testing.T) db.User {
	hashedPassword, err := util.HasPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := db.CreateUserParams{
		Username:       util.RandomOwner(),
		FullName:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	return user
}
func TestCreateUser(t *testing.T) {
	createUser(t)

}

func TestGetUser(t *testing.T) {
	user1 := createUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, 0)
	require.NotEmpty(t, user2.HashedPassword)
	// require.NotEqual(t, user1.HashedPassword, user2.HashedPassword) // Ensure password is hashed

}

func TestChangePassword(t *testing.T) {
	user := createUser(t)
	arg := db.ChangePasswordParams{
		Username:       user.Username,
		HashedPassword: "root",
	}
	user2, err := testQueries.ChangePassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.FullName, user2.FullName)
	require.Equal(t, user.Email, user2.Email)
	require.NotEqual(t, user.HashedPassword, user2.HashedPassword) // Ensure password has changed
	require.WithinDuration(t, user.CreatedAt, user2.CreatedAt, 0)
}
func TestDeleteUser(t *testing.T) {
	user := createUser(t)
	err := testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)

	// Check that the account has been deleted
	_, err = testQueries.GetUser(context.Background(), user.Username)
	require.Error(t, err)
}

func TestListUser(t *testing.T) {
	// Create 10 random accounts
	for i := 0; i < 10; i++ {
		createUser(t)
	}

	// List the accounts
	arg := db.ListUserParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListUser(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
