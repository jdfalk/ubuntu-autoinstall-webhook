// main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/cmd"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/observability"
)

func main() {
	// Initialize OpenTelemetry tracer.
	shutdown, err := observability.InitTracer("ubuntu-autoinstall-webhook")
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	// Ensure tracer shutdown on exit.
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	// Capture OS signals for graceful shutdown.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		log.Println("Received shutdown signal, exiting...")
		os.Exit(0)
	}()

	cmd.Execute()
}
