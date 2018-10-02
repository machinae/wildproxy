package main

import (
	"log"
	"net/http"
	"net/http/httputil"
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
	return &httputil.ReverseProxy{
		Director:       proxyRequest,
		ModifyResponse: proxyResponse,
	}
}
