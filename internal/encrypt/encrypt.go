// Package encrypt provides utilities for encrypting and decrypting
// sensitive values in .env files using AES-GCM with a passphrase-derived key.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const prefix = "enc:"

// IsEncrypted reports whether a value was encrypted by this package.
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, prefix)
}

// deriveKey produces a 32-byte AES key from the given passphrase using SHA-256.
func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

// Encrypt encrypts plaintext using AES-GCM and returns a prefixed base64 string.
func Encrypt(plaintext, passphrase string) (string, error) {
	if passphrase == "" {
		return "", errors.New("passphrase must not be empty")
	}
	block, err := aes.NewCipher(deriveKey(passphrase))
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return prefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a value previously encrypted by Encrypt.
// Returns an error if the value is not in the expected format.
func Decrypt(value, passphrase string) (string, error) {
	if passphrase == "" {
		return "", errors.New("passphrase must not be empty")
	}
	if !IsEncrypted(value) {
		return "", fmt.Errorf("value does not have %q prefix", prefix)
	}
	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, prefix))
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	block, err := aes.NewCipher(deriveKey(passphrase))
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm: %w", err)
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}
	return string(plaintext), nil
}

// ApplyToEnv returns a copy of env where all encrypted values are decrypted.
// Keys whose values fail to decrypt are left as-is and their keys are collected in errs.
func ApplyToEnv(env map[string]string, passphrase string) (map[string]string, map[string]error) {
	out := make(map[string]string, len(env))
	errs := map[string]error{}
	for k, v := range env {
		if IsEncrypted(v) {
			plain, err := Decrypt(v, passphrase)
			if err != nil {
				out[k] = v
				errs[k] = err
				continue
			}
			out[k] = plain
		} else {
			out[k] = v
		}
	}
	return out, errs
}
