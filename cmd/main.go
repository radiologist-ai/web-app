package main

import (
	"context"
	"fmt"
	"github.com/radiologist-ai/web-app/internal/app"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// Create a context to handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Create a WaitGroup to keep track of running goroutines
	var wg sync.WaitGroup

	// Start the HTTP server
	wg.Add(1)
	errChan := make(chan error)
	go func() {
		errChan <- app.Run(ctx, &wg)
	}()

	// Listen for termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	select {
	case err := <-errChan:
		if err != nil {
			log.Fatal(err)
		}
	case <-signalCh:
	}

	// Start the graceful shutdown process
	fmt.Println("\nGracefully shutting down HTTP server...")

	// Cancel the context to signal the HTTP server to stop
	cancel()

	// Wait for the HTTP server to finish
	wg.Wait()

	fmt.Println("Shutdown complete.")
}
