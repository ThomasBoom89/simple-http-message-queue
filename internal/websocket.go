package internal

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"runtime"
)

type WebSocket struct {
	app            *fiber.App
	registerPool   chan *couple
	unregisterPool chan *couple
	topicManager   *TopicManager
}

type couple struct {
	Connection *websocket.Conn
	Topic      Topic
}

func NewWebSocket(app *fiber.App, topicManager *TopicManager) *WebSocket {
	registerPool := make(chan *couple, 10)
	unregisterPool := make(chan *couple, 10)

	return &WebSocket{
		app:            app,
		registerPool:   registerPool,
		unregisterPool: unregisterPool,
		topicManager:   topicManager,
	}
}

func (W *WebSocket) SetupRoutes() {
	W.appendMiddleware()
	go W.sendMessagesToClients()
	go W.handleConnectionRegisterUnregister()

	W.app.Get("/:topic/ws", websocket.New(func(c *websocket.Conn) {
		topic := Topic(c.Params("topic"))
		couple := &couple{
			Connection: c,
			Topic:      topic,
		}
		// unregister connection
		defer func() {
			W.unregisterPool <- couple
			c.Close()
		}()

		// register connection
		W.registerPool <- couple
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
				W.topicManager.AddConnectionToQueue(topic, c)
			}
		}
	}, websocket.Config{RecoverHandler: GetWebsocketPanicHandler()}))
}

func (W *WebSocket) sendMessagesToClients() {
	defer RecoverGoroutine(W.sendMessagesToClients)
	for {
		runtime.Gosched()
		for _, topic := range W.topicManager.GetTopics() {
			connection, err := W.topicManager.GetNextConnection(topic)
			if err != nil {
				continue
			}
			exists := W.topicManager.ConnectionExists(topic, connection)
			if exists == false {
				continue
			}

			message, err := W.topicManager.GetNextMessage(topic)
			if err != nil {
				W.topicManager.AddConnectionToQueue(topic, connection)
				continue
			}

			err = connection.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				W.topicManager.AddConnectionToQueue(topic, connection)
				W.topicManager.AddMessage(topic, message)
			}
		}
	}
}

func (W *WebSocket) handleConnectionRegisterUnregister() {
	defer RecoverGoroutine(W.handleConnectionRegisterUnregister)
	for {
		select {
		case foo := <-W.registerPool:
			W.topicManager.AddConnection(foo.Topic, foo.Connection)
			log.Debug("connection registered")

		case foo := <-W.unregisterPool:
			W.topicManager.RemoveConnection(foo.Topic, foo.Connection)
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
