package internal

import "github.com/gofiber/fiber/v2"

type HTTP struct {
	app   *fiber.App
	queue Queue
}

func NewHTTP(app *fiber.App, queue Queue) *HTTP {
	return &HTTP{app: app, queue: queue}
}

func (H *HTTP) SetupRoutes() {
	H.app.Post("/publish", func(c *fiber.Ctx) error {
		H.queue.Enqueue(string(c.Body()))

		return c.SendStatus(fiber.StatusOK)
	})

	H.app.Get("/subscribe", func(c *fiber.Ctx) error {
		element, err := H.queue.Dequeue()
		if err != nil {
			return c.SendString("")
		}

		return c.SendString(element)
	})
}
