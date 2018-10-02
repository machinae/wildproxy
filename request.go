package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Function that modifes the request
// Proxies requests in the form of
// http://www.example.com/https://www.upstream.com/foo/bar
// TODO support http scheme?
func proxyRequest(req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/")
	targetUrl := path + req.URL.RawQuery

	u, err := url.Parse(targetUrl)
	if err != nil {
		log.Printf("Request URL error: %s\n", err)
		return
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	req.URL = u
	req.Host = u.Host

	if verbose {
		log.Printf("Proxying request to %s", req.URL)
	}
}
