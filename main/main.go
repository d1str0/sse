package main

import (
	// "encoding/json"
	"flag"
	"fmt"
	//	"github.com/d1str0/sse"
	//	"io/ioutil"
	//	"net/mail"
	"os"
	//	"strings"
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

func main() {
	flag.StringVar(&mailDir, "mail-dir", "", "directory to load mail archives from")
	flag.Parse()

	ReadAllFiles(mailDir)

	fmt.Printf("%d total files in %d different directories!\n", fileCount, dirCount)
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
		//fmt.Printf("%s/\n", stat.Name())

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
	err = f.Close()
	if err != nil {
		fmt.Printf("Error closing file %s: %s\n", stat.Name(), err.Error())
		os.Exit(1)
	}

	//fmt.Printf("%s\n", stat.Name())

	fileCount++
}

/*
	m, err := mail.ReadMessage(f)
	if err != nil {
		fmt.Printf("Error reading mail message: %#v\n", err)
		os.Exit(1)
	}

	header := m.Header
	fmt.Println("Date:", header.Get("Date"))
	fmt.Println("From:", header.Get("From"))
	fmt.Println("To:", header.Get("To"))
	fmt.Println("Subject:", header.Get("Subject"))

	body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		fmt.Printf("Error reading all of message body: %#v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s", body)

}
*/
