package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Client represents a connected websocket client
type Client struct {
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

// Message represents the structure of our chat messages
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

// Create a new hub
func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run the hub to handle client events and message broadcasting
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			log.Println("Client registered, total clients:", len(h.clients))
		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mutex.Unlock()
			log.Println("Client unregistered, total clients:", len(h.clients))
		case message := <-h.broadcast:
			h.mutex.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.Unlock()
		}
	}
}

// Run the client's read pump
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.hub.broadcast <- message
	}
}

// Run the client's write pump
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		message, ok := <-c.send
		if !ok {
			// The hub closed the channel
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("error writing message: %v", err)
			return
		}
	}
}

// Configure websocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (for development)
	},
}

// WebSocket handler
func handleWebSocket(hub *Hub, c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return err
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Start the client routines
	go client.writePump()
	go client.readPump()

	return nil
}

func main() {
	// Create a new hub and run it in a goroutine
	hub := newHub()
	go hub.run()

	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.Static("/", "public")
	e.GET("/ws", func(c echo.Context) error {
		return handleWebSocket(hub, c)
	})

	// Start server
	port := ":8080"
	fmt.Printf("Server started on http://localhost%s\n", port)
	e.Logger.Fatal(e.Start(port))
}
