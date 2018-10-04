package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Function that modifes the response
// TODO deal with 301 redirects
func proxyResponse(r *http.Response) error {
	if r == nil {
		return errors.New("Content is empty")
	}
	removeSecHeaders(r)
	setCorsHeaders(r)

	resolveRedirect(r)

	if r.Body == nil {
		return nil
	}

	// Only modify HTML responses
	if isHtml(r) {
		if err := rewriteLinks(r); err != nil {
			return err
		}
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

// Resolve absolute URL in Location header
func resolveRedirect(r *http.Response) {
	loc := r.Header.Get("Location")
	if loc == "" || r.Request == nil {
		return
	}
	rdUrl := resolveProxyURL(r.Request.URL, loc)
	r.Header.Set("Location", rdUrl)
}

func removeSecHeaders(r *http.Response) {
	// Drop CSP header for now
	r.Header.Del("Content-Security-Policy")
	r.Header.Del("Content-Security-Policy-Report-Only")

	// Disable HSTS
	r.Header.Del("Strict-Transport-Security")

	r.Header.Del("X-Frame-Options")

}

// Returns a new body with absolute urls in <a> and <script> tags changed to
// relative URLs from the proxy
// TODO if performance is too slow, try regexes
func rewriteLinks(r *http.Response) error {
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return err
	}
	r.Body.Close()
	if r.Request != nil {
		doc.Url = r.Request.URL
	}

	// Replace links href
	doc.Find("a").Each(func(i int, el *goquery.Selection) {
		href, ok := el.Attr("href")
		if !ok {
			return
		}
		el.SetAttr("href", resolveProxyURL(doc.Url, href))
	})

	// Replace script src
	doc.Find("script").Each(func(i int, el *goquery.Selection) {
		src, ok := el.Attr("src")
		if !ok {
			return
		}
		el.SetAttr("src", resolveProxyURL(doc.Url, src))
	})

	//Forms
	doc.Find("form").Each(func(i int, el *goquery.Selection) {
		act, ok := el.Attr("action")
		if !ok || act == "" {
			return
		}
		el.SetAttr("action", resolveProxyURL(doc.Url, act))
	})

	// HTML Head
	headEl := doc.Find("head")
	if headEl == nil {
		doc.Append("head")
		headEl = doc.Find("head")
	}

	// Add <base> tag set to original root so other paths are resolved
	// correctly
	if doc.Url != nil {
		baseTag := fmt.Sprintf(`<base href="%s"/>`, doc.Url)
		headEl.PrependHtml(baseTag)
	}

	// Add dummy favicon to prevent requests to favicon.ico
	faviconTag := `<link rel="icon" href="data:,">`
	headEl.AppendHtml(faviconTag)

	// replace with modified body
	html, err := doc.Html()
	if err != nil {
		return err
	}

	r.ContentLength = int64(len(html))
	r.Header.Set("Content-Length", strconv.Itoa(len(html)))
	r.Body = ioutil.NopCloser(strings.NewReader(html))
	return nil
}

// convert a URL to relative from the proxy
func resolveProxyURL(pageUrl *url.URL, rawUrl string) string {
	if pageUrl == nil || pageUrl.Host == "" {
		return rawUrl
	}
	u, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}
	// Resolve absolute URL for the page
	nu := pageUrl.ResolveReference(u)
	// Resolve again from the proxy
	nu = rootUrl.ResolveReference(&url.URL{Path: nu.String()})
	return nu.String()
}
