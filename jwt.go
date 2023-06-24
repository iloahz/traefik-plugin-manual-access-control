package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret []byte
}

// secret is a base64 encoded string
func NewJWT(secret string) (*JWT, error) {
	buf, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	// Check if the key size is valid for AES
	if len(buf) != 16 && len(buf) != 24 && len(buf) != 32 {
		return nil, fmt.Errorf("invalid key size for AES, must be 16, 24, or 32 bytes")
	}

	return &JWT{
		secret: buf,
	}, nil
}

// generate token based on seed using aes
func (j *JWT) GenerateToken(id string) string {
	// TODO make exp configurable
	claims := &jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 30 * 12 * 3).Unix(), // 3 years
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"id":  id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		panic(err)
	}
	return signedToken
}

func (j *JWT) ValidateToken(signedToken string) (jwt.Claims, error) {
	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(signedToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
