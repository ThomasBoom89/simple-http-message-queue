package main

import (
	"github.com/ThomasBoom89/simple-http-message-queue/internal"
	grpc2 "github.com/ThomasBoom89/simple-http-message-queue/internal/grpc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"google.golang.org/grpc"
	"net"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	topicManager := internal.NewTopicManager()
	storage := internal.NewStorage(topicManager)
	storage.Load()
	defer internal.SaveOnPanic(storage)
	go internal.HandleOsSignal(storage)
	wg.Add(1)
	go startGrpc(topicManager, wg)
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		},
	})

	app.Use(recover2.New(recover2.Config{EnableStackTrace: true, StackTraceHandler: internal.GetStackTraceHandler()}))
	app.Use(logger.New())
	app.Use(cors.New())

	http := internal.NewHTTP(app, topicManager)
	http.SetupRoutes()

	websocket := internal.NewWebSocket(app, topicManager)
	websocket.SetupRoutes()

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
	wg.Wait()
}

func startGrpc(topicManager *internal.TopicManager, wg *sync.WaitGroup) {
	defer wg.Done()

	listener, err := net.Listen("tcp", ":3001")
	if err != nil {
		panic(err)
	}
	server := grpc2.NewServer(topicManager)

	gsrv := grpc.NewServer()

	// Implement Server Stub from generated files
	grpc2.RegisterMessageBrokerServer(gsrv, server)

	if err := gsrv.Serve(listener); err != nil {
		panic(err)
	}
}
