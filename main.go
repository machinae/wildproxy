package main

import (
	"log"
	"time"

	flag "github.com/spf13/pflag"
)

// Flags
var (
	httpHost string

	// timeout for upstream requests
	upstreamTimeout time.Duration
	// Timeout for client requests
	clientTimeout time.Duration

	// Verbose output
	verbose bool
)

func init() {
	flag.StringVarP(&httpHost, "host", "h", ":8080", "Host to run HTTP server on")
	flag.DurationVarP(&upstreamTimeout, "upstream-timeout", "t", 60*time.Second, "Timeout for requests to upstream servers")
	flag.DurationVarP(&clientTimeout, "client-timeout", "T", 60*time.Second, "Timeout for requests from clients to this server")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

func main() {
	flag.Parse()
	log.Printf("Starting server on %s", httpHost)

	StartServer()
}
