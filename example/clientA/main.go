package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// Default address and port as start up parameters
	addr = flag.String("addr", "localhost:8080", "http service address")

	// Id generated at program started based on timestamp
	id = time.Now().Unix()
)

func main() {

	// Enstablish a connection with the web socket server
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	fmt.Printf("connecting to %s\n", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	// Set a message channel for incoming messages from server and from user input
	serverMessage := make(chan []byte)
	userMessage := make(chan []byte)

	// Messages from server will be sent to the channel serverMessage
	go ListenMessagesFromServer(conn, serverMessage)

	go ReadUserMessage(userMessage)

	for {
		select {
		case m := <-serverMessage:

			// print message only if id != this.id
			combinedID, _ := binary.Varint(m[:8])
			if combinedID != id {
				fmt.Println(string(m[8:]))
			}

		case m := <-userMessage:

			// combine message with the id
			m = combineIDWithMessage(id, m)

			// send the message to the server web socket
			err := conn.WriteMessage(websocket.TextMessage, m)
			if err != nil {
				log.Println("write:", err)
				return
			}

		}
	}
}
