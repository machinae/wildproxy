package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testPage = `
	<html>
	  <head>
	  	<link rel="stylesheet" href="/assets/style.css" />
	  </head>
	  <body>
	  	<h1>Hello, World</h1>
		<img src="https://cdn.example.com/logo.png" />
		<img src="/static/header.jpg" />
		<a href="https://www.example.com/page/1">Page 1</a>
		<a href="/page/2">Page 2</a>
		<a href="./3">Page 3</a>
	    <script src="https://cdn.example.com/scripts/script.js"></script>
		<form action="/login">
			<button type="submit">
		</form>
	  </body>
	</html>
	`

func TestProxyResponse(t *testing.T) {
	assert := assert.New(t)

	rootUrl = &url.URL{Path: "/"}

	req, _ := http.NewRequest("GET", "https://www.example.com/page/0", nil)

	resp := &http.Response{
		Request: req,
		Body:    ioutil.NopCloser(strings.NewReader(testPage)),
		Header:  make(http.Header),
	}

	err := rewriteLinks(resp)
	assert.NoError(err)

	rawBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err)
	body := string(rawBody)

	// Base element
	assert.Contains(body, `<base href="https://www.example.com/page/0"/>`)

	// images and stylesheets should not be proxied
	assert.Contains(body, `<img src="https://cdn.example.com/logo.png"/>`)

	// absolute links rewritten to proxy
	assert.Contains(body, `<a href="/https://www.example.com/page/1">`)
	assert.Contains(body, `<a href="/https://www.example.com/page/2">`)
	assert.Contains(body, `<link rel="stylesheet" href="/https://www.example.com/assets/style.css"/>`)

	// relative links resolved with page URL
	assert.Contains(body, `<a href="/https://www.example.com/page/3">`)
	assert.Contains(body, `<form action="/https://www.example.com/login">`)

	// script src also rewritten
	assert.Contains(body, `<script src="/https://cdn.example.com/scripts/script.js">`)

}

func TestParseCSS(t *testing.T) {
	assert := assert.New(t)
	rootUrl = &url.URL{Scheme: "http", Host: "localhost", Path: "/"}
	pageUrl := &url.URL{Scheme: "http", Host: "www.example.com", Path: "/"}

	style := `
	@font-face {
		font-family:'DDG_ProximaNova';
		src:url("font/ProximaNova-Sbold-webfont.eot");
		src:url(font/ProximaNova-Sbold-webfont.eot?#iefix)
		font-weight:600;
		font-style:normal
	}
	`

	r := strings.NewReader(style)
	out := rewriteStyleUrls(pageUrl, r)

	newCss, err := ioutil.ReadAll(out)
	assert.NoError(err)

	outStr := string(newCss)

	assert.Contains(outStr, `src:url("http://localhost/http://www.example.com/font/ProximaNova-Sbold-webfont.eot");`)
}
