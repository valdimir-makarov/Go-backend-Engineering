package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func proxyToUserService(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse(target)
		if err != nil {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(url)
		r.Host = url.Host
		proxy.ServeHTTP(w, r)
	}
}
