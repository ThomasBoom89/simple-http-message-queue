package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type HTTP struct {
	app          *fiber.App
	topicManager *TopicManager
}

func NewHTTP(app *fiber.App, topicManager *TopicManager) *HTTP {
	return &HTTP{app: app, topicManager: topicManager}
}

func (H *HTTP) SetupRoutes() {
	H.app.Post("/:topic/publish", func(c *fiber.Ctx) error {
		topic := Topic(c.Params("topic"))
		body := utils.CopyBytes(c.Body())
		H.topicManager.AddMessage(topic, body)

		return c.SendStatus(fiber.StatusOK)
	})

	H.app.Get("/:topic/subscribe", func(c *fiber.Ctx) error {
		topic := Topic(c.Params("topic"))
		message, err := H.topicManager.GetNextMessage(topic)
		if err != nil {
			return c.Send([]byte{})
		}

		return c.Send(message)
	})
}
