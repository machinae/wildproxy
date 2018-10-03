package main

import (
	"errors"
	"mime"
	"net/http"
)

// Function that modifes the response
func proxyResponse(r *http.Response) error {
	if r == nil {
		return errors.New("Content is empty")
	}
	removeSecHeaders(r)
	setCorsHeaders(r)

	// Only modify HTML responses
	if !isHtml(r) {
		return nil
	}
	return nil
}

// Parses content-type to determine if page is HTML
func isHtml(r *http.Response) bool {
	contentType := r.Header.Get("Content-Type")
	ct, _, err := mime.ParseMediaType(contentType)
	// TODO sniff content-type from body
	if err == nil && ct == "text/html" {
		return true
	}
	return false
}

func setCorsHeaders(r *http.Response) {
	r.Header.Set("Access-Control-Allow-Origin", "*")
	r.Header.Set("Access-Control-Allow-Credentials", "true")
	r.Header.Set("Access-Control-Max-Age", "86400")

	if r.Request != nil {
		rm := r.Request.Header.Get("Access-Control-Request-Method")
		if rm != "" {
			r.Header.Set("Access-Control-Allow-Methods", rm)
		}
		rh := r.Request.Header.Get("Access-Control-Request-Headers")
		if rh != "" {
			r.Header.Set("Access-Control-Allow-Headers", rh)
		}
	}
}

func removeSecHeaders(r *http.Response) {
	// Drop CSP header for now
	r.Header.Del("Content-Security-Policy")
	r.Header.Del("Content-Security-Policy-Report-Only")

	r.Header.Del("X-Frame-Options")

}
