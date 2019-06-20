package main // import "github.com/xperimental/goecho"

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// Version is set to the build version when building using the build script.
	Version = "unknown"

	addr          = ":8080"
	gracefulDelay = 2 * time.Second
)

func main() {
	flag.StringVar(&addr, "addr", addr, "Address and port to listen on.")
	flag.DurationVar(&gracefulDelay, "graceful-delay", gracefulDelay, "Delay between receiving a shutdown signal and starting shutdown.")
	flag.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting hostname: %s", err)
	}

	env := os.Environ()

	server := createServer(addr, Version, hostname, env)

	shutdownErrCh := make(chan error)
	go func() {
		sigCh := make(chan os.Signal)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

		<-sigCh
		signal.Reset()

		if gracefulDelay > 0 {
			log.Printf("Waiting %s for shutdown...", gracefulDelay)
			time.Sleep(gracefulDelay)
		}

		log.Println("Shutting down...")
		shutdownErrCh <- server.Shutdown(context.Background())
	}()

	log.Printf("Listening on %s\n", addr)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error during listen: %s", err)
	}

	if err := <-shutdownErrCh; err != nil {
		log.Fatalf("Error during shutdown: %s", err)
	}
}
