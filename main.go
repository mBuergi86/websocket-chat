package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

// Client represents a connected websocket client
type Client struct {
	conn *websocket.Conn
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

var (
	event = log.Info()
	fatal = log.Fatal()
)

// Initialize zerolog
func init() {
	// Set up pretty logging for development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false})
}

// Run the hub to handle client events and message broadcasting
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			event.Msgf("Client registered, total clients: %d", len(h.clients))
		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
			h.mutex.Unlock()
			event.Msgf("Client unregistered, total clients: %d", len(h.clients))
		case message := <-h.broadcast:
			h.mutex.Lock()
			for client := range h.clients {
				go func(c *Client, msg []byte) {
					err := websocket.Message.Send(c.conn, string(msg))
					if err != nil {
						fatal.Stack().Err(err).Msg("Error sending message")
						h.unregister <- c
					}
				}(client, message)
			}
			h.mutex.Unlock()
		}
	}
}

// WebSocket handler function
func handleWebSocket(hub *Hub) websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		client := &Client{ws, hub}
		hub.register <- client

		defer func() {
			hub.unregister <- client
			ws.Close()
		}()

		for {
			var msg string
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				fatal.Stack().Err(err).Msg("Error receiving message")
				break
			}
			hub.broadcast <- []byte(msg)
		}
	})
}

// WebSocket handler
func webSocketHandler(hub *Hub) echo.HandlerFunc {
	return echo.WrapHandler(handleWebSocket(hub))
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
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		Skipper:            middleware.DefaultSkipper,
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		HSTSMaxAge:         0,
	}))

	// Routes
	e.Static("/", "public")
	e.GET("/ws", webSocketHandler(hub))

	// Start server
	port := ":8080"
	fmt.Printf("Server started on http://localhost%s\n", port)
	e.Logger.Fatal(e.Start(port))
}
