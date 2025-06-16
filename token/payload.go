package token

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrExpiredToken    = errors.New("token has invalid claims: token is expired")
	ErrInvalidToken    = errors.New("invalid token")
	ErrInvalidUsername = errors.New("invalid username")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	jwt.RegisteredClaims
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        id,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	// Sử dụng Validator để validate các trường like exp, nbf, iat
	validator := jwt.NewValidator() // mặc định kiểm tra expiration & not-before
	if err := validator.Validate(p); err != nil {
		return ErrExpiredToken
	}

	if p.ID == uuid.Nil {
		return ErrInvalidToken
	}
	if p.Username == "" {
		return ErrInvalidUsername
	}
	return nil
}
