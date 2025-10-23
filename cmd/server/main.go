package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rxmy43/support-platform/internal/app"
)

func main() {
	appCtx := app.InitApp()

	go app.StartServer(appCtx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	log.Println("Server stopped gracefully")
	time.Sleep(1 * time.Second)
}
