package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
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
			dec := gob.NewDecoder(c)
			var message map[string]string
			uname := ""
			for {
				err = dec.Decode(&message)
				if err != nil {
					clients.Delete(uname)
					return
				}

				// Add client encoder to clients map on receiving first message, assuming it is the username from client
				if uname == "" {
					clients.Store(message["from"], gob.NewEncoder(c))
					uname = message["from"]
					continue
				}

				// Pass message to correct client
				go func() {

					// Gets outgoing encoder to send outgoing message to based on username
					v, ok := clients.Load(message["to"])
					if ok {
						outgoingEnc := v.(*gob.Encoder)
						outgoingEnc.Encode(message)
					} else {
						v, _ := clients.Load(message["from"])
						outgoingEnc := v.(*gob.Encoder)
						message["message"] = "User '" + message["to"] + "' not found!"
						message["from"] = message["SERVER"]
						outgoingEnc.Encode(message)
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

	// Awaits user input of EXIT and sends outgoing message to all clients signaling server closed
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')

			if text == "EXIT" {
				clients.Range(func(key, value interface{}) bool {
					outgoingEnc := value.(*gob.Encoder)
					messagemap := make(map[string]string)
					messagemap["EXIT"] = "TRUE"
					outgoingEnc.Encode(messagemap)
					return true
				})
			}
		}
	}()

	handle_connections(l, clients)
}
