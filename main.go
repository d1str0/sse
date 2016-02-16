package main

import (
	"fmt"
	"os"
)

func main() {
	key := []byte("1234567890123456")
	fmt.Printf("Key: %s\n", key)

	test1 := []byte("sometimes I fart loudly in my sleep")
	test2 := []byte("He wasn’t exactly the boogeyman, he was the guy you called to kill the boogeyman.")
	test3 := []byte("1337")
	fmt.Printf("Test:\n\t1: %s\n\t2: %s\n\t3: %s\n\n", test1, test2, test3)

	c1, _ := Encrypt(test1, key)
	c2, _ := Encrypt(test2, key)
	c3, _ := Encrypt(test3, key)
	fmt.Printf("Encrypted:\n\t1: %x\n\t2: %x\n\t3: %x\n\n", c1, c2, c3)

	h1 := HMAC(c1, key)
	h2 := HMAC(c2, key)
	h3 := HMAC(c3, key)
	fmt.Printf("HMAC with key:\n\t1: %x\n\t2: %x\n\t3: %x\n\n", h1, h2, h3)

	b1 := CheckMAC(c1, h1, key)
	b2 := CheckMAC(c2, h2, key)
	b3 := CheckMAC(c3, h3, key)
	fmt.Printf("CheckMAC:\n\t1: %t\n\t2: %t\n\t3: %t\n\n", b1, b2, b3)

	p1, _ := Decrypt(c1, key)
	p2, _ := Decrypt(c2, key)
	p3, _ := Decrypt(c3, key)
	fmt.Printf("Decrypted:\n\t1: %s\n\t2: %s\n\t3: %s\n", p1, p2, p3)

	db, err := BoltDBOpen()
	if err != nil {
		fmt.Printf("Error creating BoldDB database: %v", err)
		os.Exit(1)
	}
	err = db.Init()
	if err != nil {
		fmt.Printf("Error creating BoldDB database: %v", err)
		os.Exit(1)
	}
}
