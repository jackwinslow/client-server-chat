package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func handle_connections(source net.Listener, clients sync.Map) {
	for {
		c, err := source.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Handle incoming messages from client
		go func() {
			needUsername := true
			dec := gob.NewDecoder(c)
			var message string
			for {
				err = dec.Decode(&message)
				if err != nil {
					fmt.Println(err)
					return
				}

				// Add client encoder to clients map on receiving first message, assuming it is the username from client
				if needUsername {
					clients.Store(message, gob.NewEncoder(c))
					needUsername = false
					continue
				}

				// Pass message to correct client
				go func() {

					// Breaks down incoming message
					fields := strings.Fields(message)
					toClient := fields[0]
					fromClient := fields[1]

					// Constructs outgoing message
					outgoingMessage := fromClient + "> "
					for i := 2; i < len(fields); i++ {
						outgoingMessage = outgoingMessage + fields[i] + " "
					}
					outgoingMessage = strings.TrimSpace(outgoingMessage)

					// Gets outgoing encoder to send outgoing message to based on username
					v, ok := clients.Load(toClient)
					if (ok) {
						outgoingEnc := v.(*gob.Encoder)
						outgoingEnc.Encode(outgoingMessage)
					}
				}()
			}
		}()
	}
}

func main() {

	// Get port from input
	port := os.Args[1]

	// Initialize receiver
	address := "127.0.0.1:" + port
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()

	var clients sync.Map

	handle_connections(l, clients)
}
