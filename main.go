package main // import "github.com/xperimental/goecho"

import (
	"context"
	"flag"
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

	addr          = ":8080"
	gracefulDelay = 2 * time.Second
	allowEnv      bool
)

type tlsConfig struct {
	cert string
	key  string
}

func main() {
	var tlsc tlsConfig
	flag.StringVar(&addr, "addr", addr, "Address and port to listen on.")
	flag.StringVar(&tlsc.cert, "cert", "", "Path to TLS certificate file.")
	flag.StringVar(&tlsc.key, "key", "", "Path to TLS key file.")
	flag.DurationVar(&gracefulDelay, "graceful-delay", gracefulDelay, "Delay between receiving a shutdown signal and starting shutdown.")
	flag.BoolVar(&allowEnv, "allow-env", allowEnv, "Allow retrieval of environment variables.")
	flag.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting hostname: %s", err)
	}

	env := []string{}
	if allowEnv {
		env = os.Environ()
		sort.Strings(env)
	}

	server, unreadyFunc := createServer(addr, Version, hostname, env)

	if tlsc.cert == "" && tlsc.key == "" {
		tlsc.cert = os.Getenv("TLS_CERT")
		tlsc.key = os.Getenv("TLS_KEY")
	}

	shutdownErrCh := make(chan error)
	go func() {
		sigCh := make(chan os.Signal)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

		<-sigCh
		signal.Reset()

		// Make readiness check fail
		unreadyFunc()

		if gracefulDelay > 0 {
			log.Printf("Waiting %s for shutdown...", gracefulDelay)
			time.Sleep(gracefulDelay)
		}

		log.Println("Shutting down...")
		shutdownErrCh <- server.Shutdown(context.Background())
	}()

	log.Printf("Listening on %s\n", addr)
	if tlsc.cert != "" && tlsc.key != "" {
		err = server.ListenAndServeTLS(tlsc.cert, tlsc.key)
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
