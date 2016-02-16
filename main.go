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
	fmt.Printf("Decrypted:\n\t1: %s\n\t2: %s\n\t3: %s\n\n", p1, p2, p3)

	db, err := BoltDBOpen()
	if err != nil {
		fmt.Printf("Error creating BoldDB database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.Init()
	if err != nil {
		fmt.Printf("Error creating BoldDB database: %v", err)
		os.Exit(1)
	}
	db.Put(DOCUMENTS, h1, c1)
	if err != nil {
		fmt.Printf("Error creating BoldDB database: %v", err)
		os.Exit(1)
	}
	db.Put(DOCUMENTS, h2, c2)
	if err != nil {
		fmt.Printf("Error creating BoldDB database: %v", err)
		os.Exit(1)
	}
	db.Put(DOCUMENTS, h3, c3)
	if err != nil {
		fmt.Printf("Error creating BoldDB database: %v", err)
		os.Exit(1)
	}

	cc1, err := db.Get(DOCUMENTS, h1)
	if err != nil {
		fmt.Printf("Error creating BoldDB database: %v", err)
		os.Exit(1)
	}
	bb1 := CheckMAC(cc1, h1, key)

	cc2, err := db.Get(DOCUMENTS, h2)
	if err != nil {
		fmt.Printf("Error creating BoldDB database: %v", err)
		os.Exit(1)
	}
	bb2 := CheckMAC(cc2, h2, key)

	cc3, err := db.Get(DOCUMENTS, h3)
	if err != nil {
		fmt.Printf("Error creating BoldDB database: %v", err)
		os.Exit(1)
	}
	bb3 := CheckMAC(cc3, h3, key)
	fmt.Printf("CheckMAC (after DB):\n\t1: %t\n\t2: %t\n\t3: %t\n\n", bb1, bb2, bb3)

	err = db.Delete(DOCUMENTS, h1)
	err = db.Delete(DOCUMENTS, h2)
	err = db.Delete(DOCUMENTS, h3)
	if err != nil {
		fmt.Printf("Error creating BoldDB database: %v", err)
		os.Exit(1)
	}
}
