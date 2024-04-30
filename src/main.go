package main

import (
	"github.com/ThomasBoom89/simple-http-message-queue/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	queue := internal.NewQueue()

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

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
