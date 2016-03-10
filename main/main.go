package main

import (
	//"encoding/json"
	"fmt"
	"github.com/d1str0/sse"
	"os"
)

func main() {
	password := []byte("hunter2")
	salt, _ := sse.Salt()

	key := sse.Key(password, salt, 4096)
	fmt.Printf("Key: %x\n", key)

	test1 := []byte("sometimes I fart loudly in my sleep")
	test2 := []byte("He wasn’t exactly the boogeyman, he was the guy you called to kill the boogeyman.")
	test3 := []byte("1337")
	fmt.Printf("Test:\n\t1: %s\n\t2: %s\n\t3: %s\n\n", test1, test2, test3)

	c1, _ := sse.Encrypt(test1, key)
	c2, _ := sse.Encrypt(test2, key)
	c3, _ := sse.Encrypt(test3, key)
	fmt.Printf("Encrypted:\n\t1: %x\n\t2: %x\n\t3: %x\n\n", c1, c2, c3)

	h1 := sse.HMAC(c1, key)
	h2 := sse.HMAC(c2, key)
	h3 := sse.HMAC(c3, key)
	fmt.Printf("HMAC with key:\n\t1: %x\n\t2: %x\n\t3: %x\n\n", h1, h2, h3)

	b1 := sse.CheckMAC(c1, h1, key)
	b2 := sse.CheckMAC(c2, h2, key)
	b3 := sse.CheckMAC(c3, h3, key)
	fmt.Printf("CheckMAC:\n\t1: %t\n\t2: %t\n\t3: %t\n\n", b1, b2, b3)

	p1, _ := sse.Decrypt(c1, key)
	p2, _ := sse.Decrypt(c2, key)
	p3, _ := sse.Decrypt(c3, key)
	fmt.Printf("Decrypted:\n\t1: %s\n\t2: %s\n\t3: %s\n\n", p1, p2, p3)

	var db sse.DBConn
	db, err := sse.BoltDBOpen()
	if err != nil {
		fmt.Printf("Error creating BoltDB database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	//prev := db.Conn.Stats()

	err = db.Init()
	if err != nil {
		fmt.Printf("Error on db.Init(): %v", err)
		os.Exit(1)
	}
	err = db.Put(sse.DOCUMENTS, h1, c1)
	if err != nil {
		fmt.Printf("Error putting document1: %v", err)
		os.Exit(1)
	}
	db.Put(sse.DOCUMENTS, h2, c2)
	db.Put(sse.DOCUMENTS, h3, c3)

	cc1, err := db.Get(sse.DOCUMENTS, h1)
	if err != nil {
		fmt.Printf("Error getting h1: %v", err)
		os.Exit(1)
	}
	bb1 := sse.CheckMAC(cc1, h1, key)

	cc2, _ := db.Get(sse.DOCUMENTS, h2)
	bb2 := sse.CheckMAC(cc2, h2, key)

	cc3, _ := db.Get(sse.DOCUMENTS, h3)
	bb3 := sse.CheckMAC(cc3, h3, key)
	fmt.Printf("CheckMAC (after DB):\n\t1: %t\n\t2: %t\n\t3: %t\n\n", bb1, bb2, bb3)

	err = db.Delete(sse.DOCUMENTS, h1)
	if err != nil {
		fmt.Printf("Error deleting h1: %v", err)
		os.Exit(1)
	}
	db.Delete(sse.DOCUMENTS, h2)
	db.Delete(sse.DOCUMENTS, h3)

	c := sse.NewClient(db)
	c.SetKey("hunter2", "farts", 4096)
	err = c.Put("test1", test1)
	if err != nil {
		fmt.Printf("Error putting test1: %v", err)
		os.Exit(1)
	}
	c.Put("test2", test2)
	c.Put("test3", test3)

	doc1, err := c.Get("test1")
	if err != nil {
		fmt.Printf("Error getting test1: %v", err)
		os.Exit(1)
	}
	doc2, _ := c.Get("test2")
	doc3, _ := c.Get("test3")

	fmt.Printf("doc1: %s\ndoc2: %s\ndoc3: %s\n", doc1, doc2, doc3)

	err = c.AddDocToKeyword("farts", "test1")
	if err != nil {
		fmt.Printf("Error adding farts to test1: %v", err)
		os.Exit(1)
	}
	c.AddDocToKeyword("farts", "test2")
	c.AddDocToKeyword("farts", "test3")

	ids, err := c.Search("farts")
	if err != nil {
		fmt.Printf("Error searching for farts: %v", err)
		os.Exit(1)
	}
	for _, s := range ids {
		fmt.Println(s)
	}

	// Grab the current stats and diff them.
	//stats := db.Conn.Stats()
	//diff := stats.Sub(&prev)

	// Encode stats to JSON and print to STDERR.
	//json.NewEncoder(os.Stdout).Encode(diff)
}
