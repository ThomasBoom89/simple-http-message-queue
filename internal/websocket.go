package internal

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

type WebSocket struct {
	app            *fiber.App
	queue          Queue
	connectionPool map[*websocket.Conn]bool
	registerPool   chan *websocket.Conn
	unregisterPool chan *websocket.Conn
}

func NewWebSocket(app *fiber.App, queue Queue) *WebSocket {
	connectionPool := make(map[*websocket.Conn]bool)
	registerPool := make(chan *websocket.Conn, 10)
	unregisterPool := make(chan *websocket.Conn, 10)

	return &WebSocket{
		app:            app,
		queue:          queue,
		connectionPool: connectionPool,
		registerPool:   registerPool,
		unregisterPool: unregisterPool,
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
			_, _, err := c.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
					log.Debug(err)
				}
				return
			}
		}
	}))

}

func (W *WebSocket) sendMessagesToClients() {
	for {
		time.Sleep(1 * time.Second)
		if len(W.connectionPool) == 0 {
			continue
		}
		message, err := W.queue.Dequeue()
		if err != nil {
			continue
		}

		for connection, _ := range W.connectionPool {
			connection.WriteMessage(websocket.TextMessage, []byte(message))
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
