package internal

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func GetWebsocketPanicHandler() func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		if err := recover(); err != nil {
			HandlePanic(err)
		}
	}
}

func RecoverGoroutine(function func()) {
	if err := recover(); err != nil {
		HandlePanic(err)
		go function()
	}
}

func GetStackTraceHandler() func(c *fiber.Ctx, err interface{}) {
	return func(c *fiber.Ctx, err interface{}) {
		HandlePanic(err)
	}
}

func HandlePanic(err any) {
	log.Debug(err)
}
