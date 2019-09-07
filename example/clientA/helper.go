package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func combineIDWithMessage(id int64, message []byte) []byte {

	// make a buffer of 64 bytes
	buffer := make([]byte, 64)

	// add an int64 id into a []byte buffer
	binary.PutVarint(buffer, id)

	// concatenate message to buffer
	n := copy(buffer[8:], message)

	// return the combined message and truncate
	return buffer[:n+8]

}

// ListenMessagesFromServer for incoming messages and send it to the channel
func ListenMessagesFromServer(conn *websocket.Conn, msg chan<- []byte) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		msg <- message
	}
}

// ReadUserMessage use a bufferio to read user input from the console
func ReadUserMessage(msg chan<- []byte) {

	// Prepare a Reader for input message from the console
	reader := bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadBytes('\n')
		msg <- message
	}
}
