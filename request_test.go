package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyRequest(t *testing.T) {
	assert := assert.New(t)
	testUrls := []string{
		"www.example.com/foo",
		"www.example.com/foo?q=1",
		"http://www.example.com/foo",
		"https://www.example.com/foo",
		"https://www.example.com/foo?q=1",
	}

	for _, u := range testUrls {
		reqUrl := "http://localhost:11080/" + u
		req, err := http.NewRequest("GET", reqUrl, nil)
		assert.NoError(err)

		proxyRequest(req)

		assert.Equal("www.example.com", req.URL.Host, u)
		assert.Equal("www.example.com", req.Host, u)
		assert.Equal("/foo", req.URL.Path, u)
	}

}
