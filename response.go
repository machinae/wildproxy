package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Javascript to inject into the page
// Monkey patches XHR to proxy URLs
// Source: https://github.com/Rob--W/cors-anywhere
// Source: https://stackoverflow.com/questions/5202296/add-a-hook-to-all-ajax-requests-on-a-page
// TODO look into relative urls in scripts like rel2abs
var injectScript = `
(function() {
    var origin = window.location.protocol + '//' + window.location.host;
    var open = XMLHttpRequest.prototype.open;
    XMLHttpRequest.prototype.open = function() {
        var args = [].slice.call(arguments);
        var targetOrigin = /^https?:\/\/([^\/]+)/i.exec(args[1]);
        if (targetOrigin && targetOrigin[0].toLowerCase() !== origin) {
            args[1] = origin + '/' + args[1];
        }
        return open.apply(this, args);
    };
})();
`

//regexp to match url() in stylesheets
var cssRegex = regexp.MustCompile(`url\(['"]?(.+?)['"]?\)`)

// Function that modifes the response
// TODO deal with 301 redirects
func proxyResponse(r *http.Response) error {
	if r == nil {
		return errors.New("Content is empty")
	}
	removeSecHeaders(r)
	setCorsHeaders(r)

	resolveRedirect(r)

	if len(r.Header.Get("Set-Cookie")) > 0 {
		resolveCookies(r)
	}

	if r.Body == nil {
		return nil
	}

	// Modify URLs in HTML responses
	if isContentType("text/html", r) {
		defer r.Body.Close()
		if err := rewriteLinks(r); err != nil {
			return err
		}
	} else if isContentType("text/css", r) {
		defer r.Body.Close()
		br := rewriteStyleUrls(r.Request.URL, r.Body)
		r.Body = ioutil.NopCloser(br)
	}

	return nil
}

// Parses content-type header, for example text/html
func isContentType(contentType string, r *http.Response) bool {
	ctValue := r.Header.Get("Content-Type")
	ct, _, err := mime.ParseMediaType(ctValue)
	// TODO sniff content-type from body
	if err == nil && ct == contentType {
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

// Remove various browser security headers like CSP and HSTS to fix framing
// issues and improve privacy by disabling reporting
func removeSecHeaders(r *http.Response) {
	secHeaders := []string{
		"Content-Security-Policy",
		"Content-Security-Policy-Report-Only",
		"Expect-CT",
		"Public-Key-Pins",
		"Public-Key-Pins-Report-Only",
		"Strict-Transport-Security",
		"X-Content-Type-Options",
		"X-Frame-Options",
	}

	for _, h := range secHeaders {
		r.Header.Del(h)
	}
}

// Returns a new body with absolute urls in <a> and <script> tags changed to
// relative URLs from the proxy
// Does not close body
// TODO if performance is too slow, try regexes
func rewriteLinks(r *http.Response) error {
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return err
	}
	if r.Request != nil {
		doc.Url = r.Request.URL
	}

	// Replace links and styles href
	doc.Find("a,link").Each(func(i int, el *goquery.Selection) {
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

	// Inline style tags
	doc.Find("style").Each(func(i int, el *goquery.Selection) {
		css := strings.NewReader(el.Text())
		r := rewriteStyleUrls(doc.Url, css)
		newCss, err := ioutil.ReadAll(r)
		if err == nil && len(newCss) > 0 {
			el.SetText(string(newCss))
		}
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

	// Inject script
	scriptTag := fmt.Sprintf("<script>%s</script>", injectScript)
	headEl.AppendHtml(scriptTag)

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

// Resolve url('..') in stylesheets
func rewriteStyleUrls(baseUrl *url.URL, r io.Reader) io.Reader {
	css, err := ioutil.ReadAll(r)
	if err != nil {
		return bytes.NewReader(css)
	}
	newCss := ReplaceAllStringSubmatchFunc(cssRegex, string(css), func(matches []string) string {
		if len(matches) < 2 {
			return ""
		}
		urlString := matches[1]
		resolvedUrl := resolveProxyURL(baseUrl, urlString)
		return strings.Replace(matches[0], matches[1], resolvedUrl, 1)
	})
	return strings.NewReader(newCss)
}

// convert a URL to relative from the proxy
func resolveProxyURL(pageUrl *url.URL, rawUrl string) string {
	if pageUrl == nil || pageUrl.Host == "" {
		return rawUrl
	}
	// data urls don't need to be resolved
	if strings.HasPrefix(rawUrl, "data:") {
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

// Rewrites Set-Cookie header to scope cookies to the proxied path
func resolveCookies(r *http.Response) {
	rootDomain, _, _ := net.SplitHostPort(rootUrl.Host)
	var cookies []*http.Cookie
	if r == nil || r.Request == nil || r.Request.URL == nil {
		return
	}
	// namespace path
	origin := r.Request.URL.Scheme + "://" + r.Request.URL.Host
	for _, cookie := range r.Cookies() {
		cookie.Domain = rootDomain
		cookie.Path = "/" + origin + cookie.Path
		// Drop secure and httponly from cookies
		cookie.Secure = false
		cookie.HttpOnly = false
		cookies = append(cookies, cookie)
	}

	// clear old set-cookie headers and set new ones
	r.Header.Del("Set-Cookie")
	// matches behavior of http.SetCookie
	for _, cookie := range cookies {
		if v := cookie.String(); v != "" {
			r.Header.Add("Set-Cookie", v)
		}
	}
}

// http://elliot.land/post/go-replace-string-with-regular-expression-callback
func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}
