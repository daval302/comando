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

	// zeros 6 bytes left
	uuid &= 0x000000FF

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
			uuidMessage, err := strconv.ParseInt(string(m[:8]), 16, 64)
			if err != nil {
				log.Fatal("Bad conversion incoming messages")
			}

			fmt.Println(DecodeUUID(uuidMessage) + string(m[8:]))

		default:
			fmt.Print("send: ")
			message, _ := reader.ReadString('\n')

			// Encode message ID to UUID
			idMessage := strconv.FormatInt(EncodeUUID(conf.ID), 16)

			// Append the uuid to the the message
			message = idMessage + ": " + message

			// send the message to the server web socket
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("write:", err)
				return
			}

		}
	}
}
