package main

import (
	"flag"
	"fmt"
	"github.com/d1str0/sse"
	"io/ioutil"
	"net/mail"
	"os"
	"strconv"
)

var mailDir string // Directory to load

var fileCount int
var dirCount int

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
	flag.StringVar(&mailDir, "mail-dir", "", "directory to load mail archives from")
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

	ReadAllFiles(mailDir)

	fmt.Printf("%d total files in %d different directories!\n", fileCount, dirCount)
	ids, err := c.Search("justin.boyd@enron.com")
	if err != nil {
		fmt.Printf("Error searching database: %v", err)
		os.Exit(1)
	}
	fmt.Println("Found these IDs:")
	for _, id := range ids {
		fmt.Println(id)
	}

}

func ReadAllFiles(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", filename, err.Error())
		os.Exit(1)
	}

	stat, err := f.Stat()
	if err != nil {
		fmt.Printf("Error opening file stat %s: %s\n", filename, err.Error())
		os.Exit(1)
	}

	if stat.IsDir() {
		files, err := f.Readdir(0)
		if err != nil {
			fmt.Printf("Error reading directory %s: %s\n", stat.Name(), err.Error())
			os.Exit(1)
		}
		fmt.Printf("%s/\n", stat.Name())

		err = f.Close()
		if err != nil {
			fmt.Printf("Error closing directory %s: %s\n", stat.Name(), err.Error())
			os.Exit(1)
		}

		for _, file := range files {

			ReadAllFiles(fmt.Sprintf("%s/%s", filename, file.Name()))
		}

		dirCount++
		return
	}

	ParseMail(f)

	err = f.Close()
	if err != nil {
		fmt.Printf("Error closing file %s: %s\n", stat.Name(), err.Error())
		os.Exit(1)
	}

	fileCount++
}

func ParseMail(f *os.File) {
	m, err := mail.ReadMessage(f)
	if err != nil {
		fmt.Printf("Error reading mail message: %v\n", err)
		stat, err := f.Stat()
		if err != nil {
			fmt.Printf("Error reading mail message: %v, %s\n", err, stat.Name())
		}
		return
	}

	/*
		stat, err := f.Stat()
		if err != nil {
			fmt.Printf("Error reading mail message: %v\n", err)
			os.Exit(1)
		}
	*/

	header := m.Header
	tags := []string{}

	addrs, err := mail.ParseAddressList(header.Get("From"))
	if err != nil {
		//fmt.Printf("Error reading mail message: %v, %s\n", err, stat.Name())
	} else {
		for _, addr := range addrs {
			tags = append(tags, addr.Address)
		}
	}

	addrs, err = mail.ParseAddressList(header.Get("To"))
	if err != nil {
		//fmt.Printf("Error reading mail message: %v, %s\n", err, stat.Name())
	} else {
		for _, addr := range addrs {
			tags = append(tags, addr.Address)
		}
	}

	body, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Error reading all of mail file: %v\n", err)
		os.Exit(1)
	}

	id := strconv.Itoa(fileCount)
	err = c.Put(id, body)
	if err != nil {
		fmt.Printf("Error putting mail file in DB: %v\n", err)
		os.Exit(1)
	}

	for _, tag := range tags {
		err = c.AddDocToKeyword(tag, id)
		if err != nil {
			fmt.Printf("Error adding keyword to doc: %v\n", err)
			os.Exit(1)
		}
	}

	//fmt.Printf("Added document %d\n", fileCount)

}
