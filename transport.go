package main

import (
	"errors"
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"
)

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
