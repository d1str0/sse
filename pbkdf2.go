package main

import (
	"crypto/aes"
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
)

func Key(password, salt []byte, iter int) []byte {
	return sha256.Key(password, salt, iter, aes.BlockSize, sha256.New)
}
