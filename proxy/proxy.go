package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/sirupsen/logrus"
)

// Proxy ...
type Proxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

// NewProxy ...
func NewProxy(dockerURL string, tlsconfig *tls.Config) (*Proxy, error) {
	var p *httputil.ReverseProxy

	u, err := url.Parse(dockerURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing dockerURL: %s", err)
	}

	if u.Scheme == "unix" {
		p = &httputil.ReverseProxy{
			Director:  func(_ *http.Request) {},
			Transport: UNIXTransport{u.Path},
		}
	} else if u.Scheme == "tcp" {
		if u.Port() == "2376" {
			// Docker API HTTPS Endpoint
			u.Scheme = "https"
		} else {
			u.Scheme = "http"
		}

		p = httputil.NewSingleHostReverseProxy(u)
		p.Transport = &http.Transport{
			TLSClientConfig: tlsconfig,
		}
	} else {
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	return &Proxy{target: u, proxy: p}, nil
}

// Handler ...
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

// UNIXTransport ...
type UNIXTransport struct {
	path string
}

// RoundTrip ...
func (u UNIXTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	// TODO: Implement streaming (Websocket) support for attach

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", u.path)
			},
		},
	}

	method := req.Method
	uri := fmt.Sprintf("http://%s/%s?%s", req.Host, req.URL.Path, req.URL.RawQuery)
	body := req.Body

	r, err := http.NewRequest(method, uri, body)
	if err != nil {
		log.Errorf("error creating request: %s", err)
		return
	}

	// Copy headers
	r.Header = req.Header

	return httpc.Do(r)
}
