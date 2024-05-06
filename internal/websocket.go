package internal

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

type WebSocket struct {
	app             *fiber.App
	messageQueue    Queue[string]
	connectionQueue Queue[*websocket.Conn]
	connectionPool  map[*websocket.Conn]bool
	registerPool    chan *websocket.Conn
	unregisterPool  chan *websocket.Conn
}

func NewWebSocket(app *fiber.App, messageQueue Queue[string], connectionQueue Queue[*websocket.Conn]) *WebSocket {
	connectionPool := make(map[*websocket.Conn]bool)
	registerPool := make(chan *websocket.Conn, 10)
	unregisterPool := make(chan *websocket.Conn, 10)

	return &WebSocket{
		app:             app,
		messageQueue:    messageQueue,
		connectionQueue: connectionQueue,
		connectionPool:  connectionPool,
		registerPool:    registerPool,
		unregisterPool:  unregisterPool,
	}
}

func (W *WebSocket) SetupRoutes() {
	W.appendMiddleware()
	go W.sendMessagesToClients()
	go W.handleConnectionRegisterUnregister()

	W.app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// unregister connection
		defer func() {
			W.unregisterPool <- c
			c.Close()
		}()

		// register connection
		W.registerPool <- c
		// read messages
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
					log.Debug(err)
				}
				return
			}
			if string(message) == "next" {
				W.connectionQueue.Enqueue(c)
			}
		}
	}))

}

func (W *WebSocket) sendMessagesToClients() {
	for {
		time.Sleep(10 * time.Millisecond)
		if len(W.connectionPool) == 0 || W.connectionQueue.IsEmpty() || W.messageQueue.IsEmpty() {
			continue
		}
		connection, err := W.connectionQueue.Dequeue()
		if _, ok := W.connectionPool[connection]; ok == false || err != nil {
			continue
		}

		message, err := W.messageQueue.Dequeue()
		if err != nil {
			W.connectionQueue.Enqueue(connection)
			continue
		}

		err = connection.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			W.connectionQueue.Enqueue(connection)
			W.messageQueue.Enqueue(message)
			return
		}
	}
}

func (W *WebSocket) handleConnectionRegisterUnregister() {
	for {
		select {
		case connection := <-W.registerPool:
			W.connectionPool[connection] = true
			log.Debug("connection registered")

		case connection := <-W.unregisterPool:
			delete(W.connectionPool, connection)
			log.Debug("connection unregistered")
		}
	}
}

func (W *WebSocket) appendMiddleware() {
	W.app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
}
