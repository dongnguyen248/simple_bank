package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HasPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	err = CheckPassword(password, hashedPassword)
	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

}
