package sse

import (
	"bytes"
	"encoding/hex"
	"testing"
)

var key = []byte("some key I guess")

var messages = []string{
	"farts",
	"the quick brown fox jumped over the lazy dog",
	"this string is even longer than all of the previous ones combined because we might want to test it at this length or not idk",
	"00000000000",
}

var hashes = []string{
	"ffc3e3b2dd78a043f7b9113ff922c6a6f43d8c63b0b04d21fedcbaaa1ba936f7",
	"f107d878967be95c02a4578a2bfcd112e9f8bbd1e747e4a3e74b4f5cb80cc02b",
	"4bfa451d577867e69a4d444755bae578994c443935fe409d4e0ae67c99f827ab",
	"5693bc644fe7914dbf913152be47b06ff7a923c941ccf6753dbb8f01cecc641e",
}

func TestHMAC(t *testing.T) {
	for i, m := range messages {
		h := HMAC([]byte(m), key)
		expected, _ := hex.DecodeString(hashes[i])
		if !bytes.Equal(h, expected) {
			t.Errorf("Hash not equal to expected output.\n")
		}
	}
}

func TestCheckMAC(t *testing.T) {
	for i, h := range hashes {
		expected, _ := hex.DecodeString(h)
		b := CheckMAC([]byte(messages[i]), expected, key)
		if !b {
			t.Errorf("CheckMAC returned false.\n")
		}
	}
}
