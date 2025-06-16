package token

import "time"

type Maker interface {
	// CreateToken creates a new token for a specific username and duration.
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken checks if the token is valid and returns the username if it is.
	VerifyToken(token string) (*Payload, error)
	// RefreshToken generates a new token for a specific username and duration.
	// RefreshToken(token string, duration int64) (string, error)
	// GetTokenType returns the type of the token.
	// GetTokenType() string
	// // GetTokenDuration returns the duration of the token.
	// GetTokenDuration() int64
	// GetTokenIssuer returns the issuer of the token.
	// GetTokenIssuer() string
	// GetTokenAudience returns the audience of the token.

	// GetTokenAudience() string
	// GetTokenSubject returns the subject of the token.
	// GetTokenSubject() string
}
