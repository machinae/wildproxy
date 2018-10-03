package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyResponse(t *testing.T) {
	assert := assert.New(t)

	testPage := `
	<html>
	  <head>
	  </head>
	  <body>
	  	<h1>Hello, World</h1>
		<img src="https://cdn.example.com/logo.png" />
		<a href="https://www.example.com/page/1">Page 1</a>
		<a href="/page/2">Page 2</a>
		<a href="./3">Page 3</a>
	    <script src="https://cdn.example.com/scripts/script.js"></script>
	  </body>
	</html>
	`

	req, _ := http.NewRequest("GET", "https://www.example.com/page/0", nil)

	resp := &http.Response{
		Request: req,
		Body:    ioutil.NopCloser(strings.NewReader(testPage)),
	}

	err := rewriteLinks(resp)
	assert.NoError(err)

	rawBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err)
	body := string(rawBody)

	// images should not be proxied
	assert.Contains(body, `<img src="https://cdn.example.com/logo.png"/>`)

	// absolute links rewritten to proxy
	assert.Contains(body, `<a href="/https://www.example.com/page/1">`)
	assert.Contains(body, `<a href="/https://www.example.com/page/2">`)

	// relative links resolved with page URL
	assert.Contains(body, `<a href="/https://www.example.com/page/3">`)

	// script src also rewritten
	assert.Contains(body, `<script src="/https://cdn.example.com/scripts/script.js">`)

}
