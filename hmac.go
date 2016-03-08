package sse

import (
	"crypto/hmac"
	"crypto/sha256"
)

const HMACSize = sha256.Size

// HMAC will compute the MAC with the given message and given key.
func HMAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func DeriveKeys(key []byte) (aes_key, mac_key []byte) {
	aes_info := append([]byte("AES-Key"), []byte(One)[:]...)
	mac_info := append([]byte("MAC-Key"), []byte(One)[:]...)
	aes_key = HMAC(aes_info, key)
	mac_key = HMAC(mac_info, key)
	return
}

// CheckMAC reports whether messageMAC is a valid HMAC tag for message.
func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

// https://golang.org/pkg/crypto/hmac/
