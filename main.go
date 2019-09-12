package main // import "github.com/xperimental/goecho"

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

var (
	// Version is set to the build version when building using the build script.
	Version = "unknown"
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("Error in configuration: %s", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting hostname: %s", err)
	}

	env := []string{}
	if cfg.AllowEnv {
		env = os.Environ()
		sort.Strings(env)
	}

	server, unreadyFunc := createServer(cfg.Addr, Version, hostname, env)

	shutdownErrCh := make(chan error)
	go func() {
		sigCh := make(chan os.Signal)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

		<-sigCh
		signal.Reset()

		// Make readiness check fail
		unreadyFunc()

		if cfg.GracefulDelay > 0 {
			log.Printf("Waiting %s for shutdown...", cfg.GracefulDelay)
			time.Sleep(cfg.GracefulDelay)
		}

		log.Println("Shutting down...")
		shutdownErrCh <- server.Shutdown(context.Background())
	}()

	log.Printf("Listening on %s\n", cfg.Addr)
	if cfg.TLS.CertFile != "" && cfg.TLS.KeyFile != "" {
		err = server.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
	} else {
		err = server.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error during listen: %s", err)
	}

	if err := <-shutdownErrCh; err != nil {
		log.Fatalf("Error during shutdown: %s", err)
	}
}
