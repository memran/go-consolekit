package console

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

func HashSHA256(value string) string {
	h := sha256.Sum256([]byte(value))
	return hex.EncodeToString(h[:])
}

func HashSHA256Short(value string) string {
	return HashSHA256(value)[:16]
}

func VerifySHA256(value, hash string) bool {
	return HashSHA256(value) == hash
}

func RandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func RandomString(length int) (string, error) {
	b, err := RandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b)[:length], nil
}

func Encrypt(plaintext, key string) (string, error) {
	keyBytes := sha256Hash([]byte(key))
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext, key string) (string, error) {
	keyBytes := sha256Hash([]byte(key))
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func sha256Hash(key []byte) []byte {
	h := sha256.Sum256(key)
	return h[:]
}
