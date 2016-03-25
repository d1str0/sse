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
	salt:     "So Salty",
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

	c := Jaccard(set1, set2)
	fmt.Printf("Jaccard coefficient is %f\n", c)

}

func Jaccard(set1, set2 []string) float64 {
	m := make(map[string]bool)
	for _, s := range set1 {
		m[s] = true
	}

	var intCount int
	for _, s2 := range set2 {
		if m[s2] {
			intCount++
		}
	}

	return float64(intCount) / (float64(len(set1)) + float64(len(set2)) - float64(intCount))

}

func GetIDsForUser(search string) []string {
	ids, err := c.Search(search)
	if err != nil {
		fmt.Printf("Error searching database: %v", err)
		os.Exit(1)
	}
	return ids
}
