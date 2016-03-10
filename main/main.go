package main

import (
	"fmt"
	//	"github.com/d1str0/sse"
	"io/ioutil"
	"net/mail"
	"os"
	//	"strings"
)

const mailFile = "1.txt"

func main() {
	f, err := os.Open(mailFile)
	if err != nil {
		fmt.Printf("Error opening file: %#v\n", err)
		os.Exit(1)
	}

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
