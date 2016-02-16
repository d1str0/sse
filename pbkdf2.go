package sse

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
	"io"
)

func Key(password, salt []byte, iter int) []byte {
	return pbkdf2.Key(password, salt, iter, aes.BlockSize, sha256.New)
}

func Salt() ([]byte, error) {
	salt := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}
