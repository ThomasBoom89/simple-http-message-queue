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

	queue := internal.NewLinkedListQueue()

	http := internal.NewHTTP(app, queue)
	http.SetupRoutes()

	websocket := internal.NewWebSocket(app, queue)
	websocket.SetupRoutes()

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
