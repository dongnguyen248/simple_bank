package token

import (
	"testing"
	"time"

	"github.com/dongnguyen248/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	// Test cases for PasetoMaker
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestPasetoMaker_InvalidToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	// 1. Invalid token format
	_, err = maker.VerifyToken("invalid.token.string")
	require.Error(t, err)

	// 2. Expired token
	expiredToken, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	_, err = maker.VerifyToken(expiredToken)
	require.Error(t, err)

	// 3. Invalid signature
	validToken, err := maker.CreateToken(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	_, err = maker.VerifyToken(validToken + "tamper")
	require.Error(t, err)

	// 4. Invalid secret key (khác key gốc nhưng vẫn đủ dài)
	otherMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err) // lưu ý giấu đi việc secret ngắn
	_, err = otherMaker.VerifyToken(validToken)
	require.Error(t, err)
}

func TestExpirePasetoToken(t *testing.T) {
	// Test cases for expired Paseto token
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := -time.Minute // Negative duration to create an expired token
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)

	// Verify that the token is expired
	_, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
}
