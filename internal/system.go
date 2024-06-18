package internal

import (
	"github.com/gofiber/fiber/v2/log"
	"os"
	"os/signal"
	"syscall"
)

func HandleOsSignal(storage *Storage) {
	sigchan := make(chan os.Signal, 1)
	signals := []os.Signal{os.Kill, os.Interrupt, syscall.SIGTERM}
	signal.Notify(sigchan, signals...)
	<-sigchan

	storage.Save()

	log.Debug("os signal received")

	os.Exit(0)
}
