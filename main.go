package main // import "github.com/xperimental/goecho"

import (
	"flag"
	"log"
)

// Version is set to the build version when building using the build script.
var Version = "unknown"

var addr string

func main() {
	flag.StringVar(&addr, "addr", ":8080", "Address and port to listen on.")
	flag.Parse()

	server := createServer(addr, Version)

	log.Printf("Listening on %s\n", addr)
	log.Fatal(server.ListenAndServe())
}
