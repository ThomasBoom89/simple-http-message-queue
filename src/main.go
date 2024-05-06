package main

import (
	"github.com/ThomasBoom89/simple-http-message-queue/internal"
	websocket2 "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	queue := internal.NewLinkedListQueue[string]()

	http := internal.NewHTTP(app, queue)
	http.SetupRoutes()

	queue2 := internal.NewLinkedListQueue[*websocket2.Conn]()
	websocket := internal.NewWebSocket(app, queue, queue2)
	websocket.SetupRoutes()

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
