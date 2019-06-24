package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the server.
type Client struct {

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

var (
	// Set default port for this server
	addr = flag.String("addr", ":8080", "http service address")

	// set the upgrader buffer size
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
)

func handleClients() {
	for {
		select {
		case client := <-register:
			clients[client] = true
		case client := <-unregister:
			if _, ok := clients[client]; ok {
				delete(clients, client)
				close(client.send)
			}
		case message := <-broadcast:
			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(clients, client)
				}
			}
		}
	}
}

func main() {

	flag.Parse()

	// initialize channels
	broadcast = make(chan []byte)
	register = make(chan *Client)
	unregister = make(chan *Client)
	clients = make(map[*Client]bool)

	// Goroutine to handle incoming client to the server websocket
	go handleClients()

	// Expose ednpoints for incoming connections
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := &Client{conn: conn, send: make(chan []byte, 256)}
		register <- client

		// show up the client logged with the UUID connected
		// ...
		log.Println("client connected")
	})

	log.Println("Listening for imcoming connection")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
