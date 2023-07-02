package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Exp  int64  `json:"exp"` // seconds
	Iat  int64  `json:"iat"` // seconds
	Nbf  int64  `json:"nbf"` // seconds
	ID   string `json:"id"`
	Name string `json:"name"`
	Host string `json:"host"`
	IP   string `json:"ip"`
}

// GetExpirationTime implements the Claims interface.
func (m JWTClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(m.Exp, 0)), nil
}

// GetNotBefore implements the Claims interface.
func (m JWTClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(m.Nbf, 0)), nil
}

// GetIssuedAt implements the Claims interface.
func (m JWTClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(m.Iat, 0)), nil
}

// GetAudience implements the Claims interface.
func (m JWTClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

// GetIssuer implements the Claims interface.
func (m JWTClaims) GetIssuer() (string, error) {
	return "", nil
}

// GetSubject implements the Claims interface.
func (m JWTClaims) GetSubject() (string, error) {
	return "", nil
}
