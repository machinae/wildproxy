package main

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"strings"

	log "github.com/sirupsen/logrus"
)

const proxyName = "wildproxy"

// Transport that logs outgoing requests
type LoggingTransport struct {
	http.RoundTripper
}

func (tr *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if tr.RoundTripper == nil {
		return nil, errors.New("Underlying transport is nil")
	}
	rawReq, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	log.Debug(string(rawReq))

	return tr.RoundTripper.RoundTrip(req)
}

// Transport that enforces security restrictions on requests
type SafeTransport struct {
	http.RoundTripper
}

func (tr *SafeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Prevent request loops by checking Via header
	// https://blog.cloudflare.com/preventing-malicious-request-loops/
	if strings.Contains(req.Header.Get("Via"), proxyName) {
		return nil, http.ErrAbortHandler
	}
	req.Header.Set("Via", "1.1 "+proxyName)

	return tr.RoundTripper.RoundTrip(req)
}
