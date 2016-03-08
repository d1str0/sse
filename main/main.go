package main

import (
	//"encoding/json"
	"fmt"
	"github.com/d1str0/sse"
	"os"
)

func main() {
	password := "hunter2"
	salt, _ := sse.Salt()

	test1 := []byte("sometimes I fart loudly in my sleep")
	test2 := []byte("He wasn’t exactly the boogeyman, he was the guy you called to kill the boogeyman.")
	test3 := []byte("1337")
	fmt.Printf("Test:\n\t1: %s\n\t2: %s\n\t3: %s\n\n", test1, test2, test3)

	var db sse.DBConn
	db, err := sse.BoltDBOpen()
	if err != nil {
		fmt.Printf("Error creating BoldDB database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	//prev := db.Conn.Stats()

	c, err := sse.NewClient(db)
	if err != nil {
		fmt.Printf("Error creating BoltDB database: %v", err)
		os.Exit(1)
	}
	c.SetKey(password, string(salt), 4096)

	h1 := "test1"
	err = c.Put(h1, test1)
	if err != nil {
		fmt.Printf("Error creating BoltDB database: %v", err)
		os.Exit(1)
	}

	h2 := "test2"
	err = c.Put(h2, test2)
	if err != nil {
		fmt.Printf("Error creating BoltDB database: %v", err)
		os.Exit(1)
	}

	h3 := "test3"
	err = c.Put(h3, test3)
	if err != nil {
		fmt.Printf("Error creating BoltDB database: %v", err)
		os.Exit(1)
	}

	cc1, err := c.Get(h1)
	if err != nil {
		fmt.Printf("Error creating BoltDB database: %v", err)
		os.Exit(1)
	}

	cc2, err := c.Get(h2)
	if err != nil {
		fmt.Printf("Error creating BoltDB database: %v", err)
		os.Exit(1)
	}

	cc3, err := c.Get(h3)
	if err != nil {
		fmt.Printf("Error creating BoltDB database: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%s: %s\n%s: %s\n%s: %s\n", h1, cc1, h2, cc2, h3, cc3)

	err = c.Delete(h1)
	err = c.Delete(h2)
	err = c.Delete(h3)
	if err != nil {
		fmt.Printf("Error creating BoltDB database: %v", err)
		os.Exit(1)
	}

	// Grab the current stats and diff them.
	//stats := db.Conn.Stats()
	//diff := stats.Sub(&prev)

	// Encode stats to JSON and print to STDERR.
	//json.NewEncoder(os.Stdout).Encode(diff)

}
