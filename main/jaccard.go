package main

import (
	"flag"
	"fmt"
	"github.com/d1str0/sse"
	"os"
)

var user1 string
var user2 string

type Identity struct {
	password string
	salt     string
	iter     int
}

var id = Identity{
	password: "hunter2",
	salt:     "farts",
	iter:     4096,
}

var c *sse.Client

func main() {
	flag.StringVar(&user1, "user1", "", "User 1 for Jaccard coefficient")
	flag.StringVar(&user2, "user2", "", "User 2 for Jaccard coefficient")
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

	set1 := GetIDsForUser(user1)
	set2 := GetIDsForUser(user2)

	c := sse.Jaccard(set1, set2)
	fmt.Printf("Jaccard coefficient is %f\n", c)

}

func GetIDsForUser(search string) []string {
	ids, err := c.Search(search)
	if err != nil {
		fmt.Printf("Error searching database: %v", err)
		os.Exit(1)
	}
	return ids
}
