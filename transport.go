package main

import (
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
		tr.RoundTripper = &http.Transport{}
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
	if tr.RoundTripper == nil {
		tr.RoundTripper = &http.Transport{}
	}
	// Prevent request loops by checking Via header
	// https://blog.cloudflare.com/preventing-malicious-request-loops/
	if strings.Contains(req.Header.Get("Via"), proxyName) {
		return nil, http.ErrAbortHandler
	}
	req.Header.Set("Via", "1.1 "+proxyName)

	return tr.RoundTripper.RoundTrip(req)
}

// Transport that strips proxy headers like X-Forward-*
type AnonTransport struct {
	http.RoundTripper
}

func (tr *AnonTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if tr.RoundTripper == nil {
		tr.RoundTripper = &http.Transport{}
	}

	stripHeaders := []string{
		"Forwarded",
		"X-Forwarded-For",
		"X-Forwarded-Proto",
		"X-Forwarded-Host",
		"Via",
		// Cloudflare
		"CF-Connecting-IP",
		"CF-Ipcountry",
		"CF-Visitor",
		"True-Client-IP",
	}

	for _, h := range stripHeaders {
		req.Header.Del(h)
	}
	return tr.RoundTripper.RoundTrip(req)
}
