package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// Default address
	// Edit this line later as Test requirements
	addr = flag.String("addr", "localhost:8080", "http service address")
)

// EncodeUUID encode string to int64 rapresentation
func EncodeUUID(jsonID string) int64 {
	id, err := strconv.ParseInt(jsonID, 0, 16)
	if err != nil {
		log.Fatal("Bad conversion", err)
	}

	// Get the timestamp
	timestamp := time.Now().Unix()

	// zeros last 2 bytes of timestamp
	timestamp &= 0xFFFFFF00

	// add id to last 2 bytes
	timestamp += id

	return timestamp
}

// DecodeUUID decode the id from the uuid
func DecodeUUID(uuid int64) string {

	// zeros 6 bytes left
	uuid &= 0x000000FF

	// convert to string hex
	ret := strconv.FormatInt(uuid, 16)

	return ret
}

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
	newMessage := make(chan []byte)
	inputMessage := make(chan []byte)

	// Reader for input message from the console
	reader := bufio.NewReader(os.Stdin)

	// Listen messages from server
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			newMessage <- message
		}
	}()

	// Listen for input messages
	go func() {
		for {
			// fmt.Print("send: ")
			message, _ := reader.ReadBytes('\n')
			inputMessage <- message
		}

	}()

	for {
		select {
		case m := <-newMessage:

			fmt.Println(string(m))

		case m := <-inputMessage:

			// send the message to the server web socket
			err := conn.WriteMessage(websocket.TextMessage, m)
			if err != nil {
				log.Println("write:", err)
				return
			}

		}
	}
}
