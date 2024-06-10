package main

import (
	"github.com/ThomasBoom89/simple-http-message-queue/internal"
	websocket2 "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	queue := internal.NewLinkedListQueue[string]()
	storage := internal.NewStorage(queue)
	storage.Load()
	defer internal.SaveOnPanic(storage)
	go internal.HandleOsSignal(storage)

	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		},
	})

	app.Use(recover2.New(recover2.Config{EnableStackTrace: true, StackTraceHandler: internal.GetStackTraceHandler()}))
	app.Use(logger.New())
	app.Use(cors.New())

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
