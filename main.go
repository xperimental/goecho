package main // import "github.com/xperimental/goecho"

import (
	"flag"
	"log"
	"os"
)

// Version is set to the build version when building using the build script.
var Version = "unknown"

var addr string

func main() {
	flag.StringVar(&addr, "addr", ":8080", "Address and port to listen on.")
	flag.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting hostname: %s", err)
	}

	server := createServer(addr, Version, hostname)

	log.Printf("Listening on %s\n", addr)
	log.Fatal(server.ListenAndServe())
}
