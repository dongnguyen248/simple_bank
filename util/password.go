package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HasPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("fail to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	if len(hashedPassword) == 0 {
		return fmt.Errorf("hashed password is empty")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("password mismatch: %w", err)
	}
	return nil
}
