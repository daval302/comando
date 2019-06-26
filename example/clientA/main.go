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
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// Conf holds json configurations files
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

	// convert to string hex
	ret := strconv.FormatInt(uuid, 16)

	return ret
}

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
	clientID, err := strconv.ParseInt(conf.ID, 16, 16)

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
			// get the number of seconds since 01/01/1970
			timestamp := time.Now().Unix()

			// shift left 2 bytes from the timestamp
			timestamp = timestamp << 16
			// add the 2 byte id to the right
			uuid := timestamp ^ int64(clientID)

			// Append the uuid to the the message
			message = string(uuid) + ": " + message

			// send the message to the server web socket
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("write:", err)
				return
			}

		}
	}
}
