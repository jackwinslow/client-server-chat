package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
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

	gold_chain := make(chan int)

	// Sends the server username upon connection initialization
	enc := gob.NewEncoder(c)
	messagemap := make(map[string]string)
	messagemap["from"] = username
	enc.Encode(messagemap)

	// Handles message receiving
	go func() {
		dec := gob.NewDecoder(c)
		var message map[string]string
		for {
			decerr := dec.Decode(&message)
			if decerr == nil {
				if _, ok := message["EXIT"]; ok {
					gold_chain <- 1
				}
				fmt.Println(message["from"] + "> " + message["message"])
			}
		}
	}()

	// Handles further messages from user input
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			fields := strings.Fields(text)

			// Exit client (disconnection is handled by server) on EXIT command from user input
			if strings.TrimSpace(fields[0]) == "EXIT" {
				gold_chain <- 1
			}

			messagemap := make(map[string]string)
			messagemap["to"] = fields[0]
			messagemap["from"] = username
			messagemap["message"] = strings.Join(fields[1:], " ")

			// Checks to make sure all required fields are included in user input
			if len(fields) < 2 {
				fmt.Println("Please send messages in <To> <Message> format")
				continue
			}

			encerr := enc.Encode(messagemap)
			if encerr != nil {
				log.Fatal("server connection lost")
			}

		}
	}()

	// wait for kill channel to recieve signal
	_ = <-gold_chain
	fmt.Println("Client has been disconnected")
	return
}
