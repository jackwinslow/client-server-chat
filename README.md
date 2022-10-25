# client-server-chat

## Usage
To run the chat app, first initialize an instance of the server by running `go run server.go <port>` where port is the listening port the server application will run on.
To connect to the server, one must run `go run client.go <ip> <port> <username>`, which will create an instance connected to the server running on the given ip and port. 
To send a message, type `<to> <message>` to send the given `message` to the provided user, `to`, assuming the username is also connected to the server.


## message structure
the message is sent as a map between the server and clients through the gob encoding with the following fields: `{"to": <msg recepient>, "from": <msg sender>, "message" <msg>}`

## Server Side

### client map
a thread safe map of gob encoders indexed by username. 

### handle_connections
`handle_connections(source net.Listener, clients sync.Map)` oversees the routing of messages as well as creating updates and modifications to the client map. In the main thread of the function, it listens to accept client connections. After a client connects, the main thread creates a new goroutine. The newly created goroutine acts as a decoder for incoming messages, it also creates an entry in the client map for the newly connected user. When a message is recieved and decoded, it is passed into a new goroutine that determines which user to send it to, and then sends it using their encoder located in the client map

## Client Side
the client side is fairly simple as most data processing/handling is done by the server side. The client dials to the host/port with the given username. After doing so an encoder is created, and a map with only the `from` field is sent to update the client map on the server side. The process then creates a goroutine to decode messages being recieved from the server. The main thread waits for user input and encodes valid messages to the server.


### Additional notes on design choices:

According to [this PR](https://go-review.googlesource.com/c/go/+/155742), as well as [the gob documentation](https://pkg.go.dev/encoding/gob#Encoder), gob encoders and decoders are considered thread safe. Because of this, instead of creating buffer channels to feed subprocesses, we instead used a map to store encoder structs which are passed into individual goroutines to handle messages. We used a [sync map](https://pkg.go.dev/sync#Map) to make sure the encoders were stored and accessed safely as it is the only other structure that can be used by multiple threads in a potentially dangerous way. 
