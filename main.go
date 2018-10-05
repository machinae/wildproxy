package main

import (
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	flag "github.com/spf13/pflag"
)

// Flags
var (
	httpHost string
	webRoot  string

	// timeout for upstream requests
	upstreamTimeout time.Duration
	// Timeout for client requests
	clientTimeout time.Duration

	// Verbose output
	verbose bool

	// Dump outgoing requests
	debug bool

	// Anonymous proxy mode
	anonMode bool

	// Add cors headers to responses
	corsHeaders bool

	// Remove security headers
	secHeaders bool

	// Proxy all links, not just scripts and css
	rewriteAll bool
)

var (
	rootUrl *url.URL
)

func init() {
	flag.StringVarP(&httpHost, "host", "h", ":8080", "Host to run HTTP server on")
	flag.StringVarP(&webRoot, "root", "r", "", "Web root the proxy will be available at, prepended to all URLs")
	flag.DurationVarP(&upstreamTimeout, "upstream-timeout", "t", 60*time.Second, "Timeout for requests to upstream servers")
	flag.DurationVarP(&clientTimeout, "client-timeout", "T", 60*time.Second, "Timeout for requests from clients to this server")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	flag.BoolVar(&debug, "debug", false, "Dump outoging requests to debug")
	flag.BoolVar(&anonMode, "anon", false, "Strip proxy headers like X-Forwarded-For that leak user data")
	flag.BoolVar(&corsHeaders, "cors", true, "Add CORS headers to responses")
	flag.BoolVar(&secHeaders, "csp", true, "Strip content security and frame headers from responses")
	flag.BoolVarP(&rewriteAll, "all", "a", false, "Proxy all resources, not just HTML, scripts and stylesheets")
}

func main() {
	var err error
	flag.Parse()

	rootUrl, err = url.Parse(webRoot)
	if err != nil || rootUrl.Host == "" {
		log.Fatal("Root URL specified with -r must be an absolute URL like http://proxy.example.com")
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	} else if verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	compileSelectors()

	log.Printf("Starting server on %s", httpHost)
	log.Printf("Proxying requests to %s/*", rootUrl)

	StartServer()
}
