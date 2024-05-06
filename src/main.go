package main

import (
	"github.com/ThomasBoom89/simple-http-message-queue/internal"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"time"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	queue := internal.NewLinkedListQueue()

	app.Post("/publish", func(c *fiber.Ctx) error {
		queue.Enqueue(string(c.Body()))

		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/subscribe", func(c *fiber.Ctx) error {
		element, err := queue.Dequeue()
		if err != nil {
			return c.SendString("")
		}

		return c.SendString(element)
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	connectionPool := make(map[*websocket.Conn]bool, 10)
	registerPool := make(chan *websocket.Conn, 10)
	unregisterPool := make(chan *websocket.Conn, 10)

	go func() {
		for {
			select {
			case connection := <-registerPool:
				connectionPool[connection] = true
				log.Debug("connection registered")

			case connection := <-unregisterPool:
				delete(connectionPool, connection)
				log.Debug("connection unregistered")
			}
		}
	}()

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// unregister connection
		defer func() {
			unregisterPool <- c
			c.Close()
		}()

		// register connection
		registerPool <- c
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

	// send messages to pool of connection
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if len(connectionPool) == 0 {
				continue
			}
			message, err := queue.Dequeue()
			if err != nil {
				continue
			}

			for connection, _ := range connectionPool {
				connection.WriteMessage(websocket.TextMessage, []byte(message))
			}
		}
	}()

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
