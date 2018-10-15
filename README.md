# wildproxy
Wildproxy is a wildcard HTTP reverse proxy. It can proxy all requests to any
domain while stripping CORS and other security headers and rewriting
links, including scripts and stylesheets.

Unlike other reverse proxies, wildproxy does notd require extensive manual
configuration or mapping HTTP resources - it automatically proxies ANY requests
sent through it, to any domain. The typical use case for wildproxy is dealing
with cross-domain issues for Javacript, or embedding a third-party website in an
iframe.

**Security Notice: If wildproxy is deployed on a public-facing web server, it
will act as an open HTTP proxy. That means anyone can use your server IP to proxy
requests to any URL. Only deploy wildproxy on a public address if you know what
you are doing.**

**Security Notice #2: wildproxy always runs an unsecured(http) web server.
It transparently proxies all https websites to http, effectively stripping SSL from all connections.
For any public deployment, you MUST put another server like nginx or a reverse proxy like Cloudflare in front of your wildproxy server to force all connections to https. **

## Install
Install via go get:
`go get github.com/machinae/wildproxy`

## Quick Start
After installing, just launch wildproxy:
`wildproxy`

By default, the proxy server will run on **localhost:8080**. Now, you can make
requests to http://localhost:8080/$URL_TO_PROXY. The URL can be a JSON API,
Javascript, or a full HTML file. For example, to proxy a request to reddit.com,
the full URL to visit is:

`http://localhost:8080/https://www.reddit.com`

If you inspect a request to a JSON endpoint, you'll see CORS headers are
automatically added.


`http://localhost:8080/https://httpbin.org/get`


To change the host/port wildproxy listens on, use the `-h` flag:
`wildproxy -h 127.0.0.1:12080`

The other important flag to know about is `-r`, for the root URL to use when
rewriting links. For example, if the proxy is accessible from
https://proxy.corp.com:8080/, run `wildproxy -r https://proxy.corp.com:8080`.

This will transparently rewrite urls in HTML pages and AJAX XHR requests to
prefix them with the proxy's public address.


## Usage
Configuration is done through command-line flags.

```sh
$ wildproxy --help

Usage of wildproxy:
  -a, --all                         Proxy all resources, not just HTML, scripts and stylesheets
      --anon                        Strip proxy headers like X-Forwarded-For that leak user data
  -T, --client-timeout duration     Timeout for requests from clients to this server (default 1m0s)
      --cors                        Add CORS headers to responses (default true)
      --csp                         Strip content security and frame headers from responses (default true)
      --debug                       Dump outoging requests to debug
  -h, --host string                 Host to run HTTP server on (default "localhost:8080")
  -r, --root string                 Web root the proxy will be available at, prepended to all URLs
  -s, --script string               Path to Javascript file to inject in every page (default "./wildproxy.js")
  -t, --upstream-timeout duration   Timeout for requests to upstream servers (default 1m0s)
  -v, --verbose                     Verbose output

```
