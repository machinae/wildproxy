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
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
)

var (
	//regexp to match url() in stylesheets
	cssRegex          *regexp.Regexp

	// CSS selector for attributes with href attribute to rewrite
	hrefSelector goquery.Matcher
	// CSS selector for attributes with src attribute to rewrite
	srcSelector   goquery.Matcher
	mediaSelector goquery.Matcher
	linkSelector  goquery.Matcher
	scriptSelector goquery.Matcher
	formSelector goquery.Matcher
	// Selector for inline style tags
	styleSelector     goquery.Matcher
	styleAttrSelector goquery.Matcher
	securityHeaders   = []string{
		"Content-Security-Policy",
		"Content-Security-Policy-Report-Only",
		"Expect-CT",
		"Public-Key-Pins",
		"Public-Key-Pins-Report-Only",
		"Strict-Transport-Security",
		"X-Content-Type-Options",
		"X-Frame-Options",
	}
)

func compileSelectors() {
	cssRegex = regexp.MustCompile(`url\(['"]?(.+?)['"]?\)`)

	hrefSelector = cascadia.MustCompile("a,link")

	if rewriteAll {
		srcSelector = cascadia.MustCompile("[src]")
	} else {
		srcSelector = cascadia.MustCompile("script")
	}

	mediaSelector = cascadia.MustCompile("img,video,audio")

	linkSelector = cascadia.MustCompile("link")
	scriptSelector = cascadia.MustCompile("script")
	formSelector = cascadia.MustCompile("form[action]")
	styleSelector = cascadia.MustCompile("style")
	styleAttrSelector = cascadia.MustCompile("[style]")
}

// Function that modifes the response
// TODO deal with 301 redirects
func proxyResponse(r *http.Response) error {
	if r == nil {
		return errors.New("Content is empty")
	}

	if secHeaders {
		removeSecHeaders(r)
	}

	if corsHeaders {
		setCorsHeaders(r)
	}

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

	// The 'Content-Length' header is set automatically
	r.Header.Del("Content-Length")

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
	for _, h := range securityHeaders {
		r.Header.Del(h)
	}
}

func removeSecMetaTags(head *goquery.Selection) {
	for _, h := range securityHeaders {
		securityMetaTag := fmt.Sprintf(`meta[http-equiv='%s']`, h)
		securityMetaTagMatcher := cascadia.MustCompile(securityMetaTag)

		head.FindMatcher(securityMetaTagMatcher).Each(func(i int, el *goquery.Selection) {
			el.RemoveMatcher(securityMetaTagMatcher)
		})
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
	doc.FindMatcher(hrefSelector).Each(func(i int, el *goquery.Selection) {
		href, ok := el.Attr("href")
		if !ok {
			return
		}
		el.SetAttr("href", resolveProxyURL(doc.Url, href))
	})

	// Replace script src
	doc.FindMatcher(srcSelector).Each(func(i int, el *goquery.Selection) {
		el.RemoveAttr("integrity")
		el.RemoveAttr("crossorigin")

		src, ok := el.Attr("src")
		if !ok {
			return
		}
		el.SetAttr("src", resolveProxyURL(doc.Url, src))
	})

	// Replace media
	doc.FindMatcher(mediaSelector).Each(func(i int, el *goquery.Selection) {
		src, ok := el.Attr("src")
		if !ok {
			return
		}
		el.SetAttr("src", absoluteURL(doc.Url, src))
	})

	// Links
	doc.FindMatcher(linkSelector).Each(func(i int, el *goquery.Selection) {
		el.RemoveAttr("integrity")
		el.RemoveAttr("crossorigin")
	})

	//Forms
	doc.FindMatcher(formSelector).Each(func(i int, el *goquery.Selection) {
		act, ok := el.Attr("action")
		if !ok || act == "" {
			return
		}
		el.SetAttr("action", resolveProxyURL(doc.Url, act))
	})

	// Inline style tags
	doc.FindMatcher(styleSelector).Each(func(i int, el *goquery.Selection) {
		css := strings.NewReader(el.Text())
		r := rewriteStyleUrls(doc.Url, css)
		newCss, err := ioutil.ReadAll(r)
		if err == nil && len(newCss) > 0 {
			el.SetText(string(newCss))
		}
	})

	doc.FindMatcher(styleAttrSelector).Each(func(i int, el *goquery.Selection) {
		css, ok := el.Attr("style")
		if !ok {
			return
		}
		r := rewriteStyleUrls(doc.Url, strings.NewReader(css))
		newCSS, err := ioutil.ReadAll(r)
		if err == nil && len(newCSS) > 0 {
			el.SetAttr("style", string(newCSS))
		}
	})

	// HTML Head
	headEl := doc.Find("head")
	if headEl == nil {
		doc.Append("head")
		headEl = doc.Find("head")
	}

	if secHeaders {
		removeSecMetaTags(headEl)
	}

	targetScriptUrls := "window.targetScriptUrls = ["
	targetAsyncScriptUrls := "window.targetAsyncScriptUrls = ["
	// Remove external scripts
	doc.FindMatcher(scriptSelector).Each(func(i int, el *goquery.Selection) {
		srcAttr, srcAttrExist := el.Attr("src")
		_, asyncAttrExist:= el.Attr("async")
		_, deferAttrExist := el.Attr("defer")

		if !srcAttrExist || srcAttr == "" {
			return
		}

		if asyncAttrExist || deferAttrExist {
			targetAsyncScriptUrls += fmt.Sprintf(`"%s",`, srcAttr)
		} else {
			targetScriptUrls += fmt.Sprintf(`"%s",`, srcAttr)
		}

		el.Remove()
	})

	targetScriptUrls += "];"
	targetAsyncScriptUrls += "];"

	targetScriptUrlsScript := fmt.Sprintf("<script>%s%s</script>", string(targetScriptUrls), string(targetAsyncScriptUrls))
	headEl.AppendHtml(targetScriptUrlsScript)

	// Inject script
	if scriptFile != "" {
		dat, _ := ioutil.ReadFile(scriptFile)
		scriptTag := fmt.Sprintf("<script>" + string(dat) + "</script>")
		headEl.AppendHtml(scriptTag)
	}

	// replace with modified body
	html, err := doc.Html()
	if err != nil {
		return err
	}

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

		resolvedUrl := absoluteURL(baseUrl, urlString)
		if rewriteAll {
			resolvedUrl = resolveProxyURL(baseUrl, urlString)
		}

		return strings.Replace(matches[0], matches[1], resolvedUrl, 1)
	})
	return strings.NewReader(newCss)
}

// convert a URL to relative from the proxy
func resolveProxyURL(pageUrl *url.URL, rawUrl string) string {
	if pageUrl == nil || pageUrl.Host == "" {
		return rawUrl
	}
	// data urls don't and hash urls need to be resolved
	if strings.HasPrefix(rawUrl, "data:") || strings.HasPrefix(rawUrl, "#") {
		return rawUrl
	}
	u, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}

	// Remove the ending slash for correct url resolve.
	// http://host/path/ and newPath should resolve to http://host/newPath/ , but not http://host/path/newPath
	// See https://play.golang.org/p/2JfBbaXOAMy
	pageUrl.Path = strings.TrimSuffix(pageUrl.Path, "/")

	// Resolve absolute URL for the page
	nu := pageUrl.ResolveReference(u)

	nu, err = url.Parse("./" + nu.String())
	if err != nil {
		return nu.String()
	}

	// Resolve again from the proxy
	nu = rootUrl.ResolveReference(nu)
	return nu.String()
}

func absoluteURL(pageURL *url.URL, raw string) string {
	if pageURL == nil || pageURL.Host == "" {
		return raw
	}

	if strings.HasPrefix(raw, "data:") || strings.HasPrefix(raw, "http") || strings.HasPrefix(raw, "//") {
		return raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	nu := pageURL.ResolveReference(u)

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
