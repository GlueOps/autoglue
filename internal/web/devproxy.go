package web

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func DevProxy(target string) (http.Handler, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	p := httputil.NewSingleHostReverseProxy(u)
	return p, nil
}
