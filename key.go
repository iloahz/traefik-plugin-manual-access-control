package traefikpluginmanualaccesscontrol

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Key type takes a seed and generates a token
type Key struct {
	key []byte
}

// key is a base64 encoded string
func NewKey(keyB64 string) (*Key, error) {
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, err
	}

	// Check if the key size is valid for AES
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key size for AES, must be 16, 24, or 32 bytes")
	}

	return &Key{
		key: key,
	}, nil
}

// generate token based on seed using aes
func (k *Key) GenerateToken() string {
	// Create a new AES cipher
	block, err := aes.NewCipher(k.key)
	if err != nil {
		panic(err)
	}

	// Generate a random nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	// Use the GCM (Galois/Counter Mode) for authenticated encryption
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	// Encrypt the nonce using the AES-GCM cipher
	ciphertext := aesgcm.Seal(nil, nonce, nonce, nil)

	// Combine the nonce and the ciphertext
	token := append(nonce, ciphertext...)

	// Encode the token in base64
	return base64.StdEncoding.EncodeToString(token)
}

func (k *Key) ValidateToken(token string) bool {
	// Decode the token from base64
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	// Check if the decoded token has at least 12 bytes (nonce size)
	if len(decodedToken) < 12 {
		return false
	}

	// Separate the nonce and ciphertext
	nonce := decodedToken[:12]
	ciphertext := decodedToken[12:]

	// Create a new AES cipher using the seed as the key
	block, err := aes.NewCipher(k.key)
	if err != nil {
		return false
	}

	// Create a new GCM cipher using the AES block cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return false
	}

	// Decrypt the ciphertext using the GCM cipher and nonce
	_, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return false
	}

	// If decryption is successful, the token is valid
	return true
}
