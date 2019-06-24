package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

type Conf struct {
	Name string
	ID   string
}

var (
	// DTO for the json configuration file
	conf Conf

	// Default address
	// Edit this line later as Test requirements
	addr = flag.String("addr", "localhost:8080", "http service address")
)

func main() {

	// Load configuration json file
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Load configuration file")
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)

	json.Unmarshal(byteValue, &conf)

	// get 2 bytes json as ID
	clientID := []byte(conf.ID)

	fmt.Printf("client : %#x\n", clientID)

	// Enstablish a connection with the web socket server
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	fmt.Printf("connecting to %s\n", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	// Set a message channel for incoming messages from server
	newMessage := make(chan []byte)

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

	// Listen for input message from the console and incoming messages
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case m := <-newMessage:
			// decode UUID
			// ...
			fmt.Println(string(m))

		default:
			fmt.Print("send: ")
			message, _ := reader.ReadString('\n')

			// Encode message UUID
			// ...

			// send the message to the server web socket
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("write:", err)
				return
			}

		}
	}
}
