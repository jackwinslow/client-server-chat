package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Get host address, port, and username from input
	host := os.Args[1]
	port := os.Args[2]
	username := os.Args[3]

	// Setup server connection
	c, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Sends the server username upon connection initialization
	enc := gob.NewEncoder(c)
	enc.Encode(username)

	// Handles message receiving
	go func() {
		dec := gob.NewDecoder(c)
		var message string

		for {
			dec.Decode(&message)
			fmt.Println(message)
		}
	}()

	// Handles further messages from user input
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		fields := strings.Fields(text)

		// Checks to make sure all required fields are included in user input
		if len(fields) < 3 {
			fmt.Println("Please send messages in <To> <From> <Message> format")
			continue
		}

		enc.Encode(text)
	}
}
