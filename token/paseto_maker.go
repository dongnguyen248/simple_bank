package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid symmetric key size: must be %d bytes", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// CreateToken creates a new token for a specific username and duration.
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	// Thêm `nil` cho footer
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", err
	}
	return token, nil
}
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var payload Payload
	// Thêm `nil` cho footer
	if err := maker.paseto.Decrypt(token, maker.symmetricKey, &payload, nil); err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	// Kiểm tra hết hạn
	if time.Now().After(payload.ExpiredAt) {
		return nil, ErrExpiredToken
	}
	return &payload, nil
}
