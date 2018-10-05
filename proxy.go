package main

import (
	"net"
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"
)

var (
	proxy *httputil.ReverseProxy

	srv *http.Server
)

func StartServer() {
	proxy = newProxy()

	srv = &http.Server{
		Addr:         httpHost,
		Handler:      proxy,
		ReadTimeout:  clientTimeout,
		WriteTimeout: clientTimeout,
	}
	log.Fatal(srv.ListenAndServe())
}

func newProxy() *httputil.ReverseProxy {
	var tr http.RoundTripper
	dialer := &net.Dialer{
		Timeout:   upstreamTimeout,
		KeepAlive: upstreamTimeout,
	}
	tr = &http.Transport{
		DialContext: dialer.DialContext,
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		tr = &LoggingTransport{tr}
	}
	return &httputil.ReverseProxy{
		Director:       proxyRequest,
		ModifyResponse: proxyResponse,
		Transport:      tr,
	}
}
