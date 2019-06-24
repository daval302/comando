package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Set time waiting for write, pong and ping period
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Byte encding for messages
var (
	newline = []byte{'\n'}
	space   = []byte{' '}
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

func (c *Client) readPump() {
	defer func() {
		unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// send the messsage to channel broadcast in order to be sent to other clients
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Conneciton closed.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

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

		// listen for incoming messages from client
		go client.readPump()
		// broadcast messages to the client
		go client.writePump()

		log.Println("client connected")
	})

	log.Println("Listening for imcoming connection")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
