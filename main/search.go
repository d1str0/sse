package main

import (
	"flag"
	"fmt"
	"github.com/d1str0/sse"
	"os"
)

var search string

type Identity struct {
	password string
	salt     string
	iter     int
}

var id = Identity{
	password: "hunter2",
	salt:     "So Salty",
	iter:     4096,
}

var c *sse.Client

func main() {
	flag.StringVar(&search, "search", "", "search for IDs with matching keyword")
	flag.Parse()

	var db sse.DBConn
	db, err := sse.BoltDBOpen()
	if err != nil {
		fmt.Printf("Error creating BoltDB database: %v", err)
		os.Exit(1)
	}
	c, err = sse.NewClient(db)
	if err != nil {
		fmt.Printf("Error creating new client: %v", err)
		os.Exit(1)
	}
	c.SetKey("hunter2", "farts", 4096)

	ids, err := c.Search(search)
	if err != nil {
		fmt.Printf("Error searching database: %v", err)
		os.Exit(1)
	}
	fmt.Println("Found these IDs:")
	for _, id := range ids {
		fmt.Printf("%s\n", id)
	}

}
