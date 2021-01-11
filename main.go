package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	frontendURL = "http://localhost:3000"
	backendURL  = "http://localhost:8080"
)

func main() {
	proxyMux := http.NewServeMux()
	proxyMux.HandleFunc("/", forward)

	if err := http.ListenAndServe(":4000", proxyMux); err != nil {
		log.Fatalln("Server has failed to start:", err.Error())
	}
}

func forward(res http.ResponseWriter, req *http.Request) {
	isBackend := strings.HasPrefix(req.RequestURI, "/api/") // Here we assume the only backend is "api"
	if isBackend {
		req.URL.Path = strings.SplitN(req.RequestURI, "api", 2)[1]
		reverseProxyTo(backendURL, res, req)
		return
	} else {
		reverseProxyTo(frontendURL, res, req)
		return
	}
}

func reverseProxyTo(urlTo string, res http.ResponseWriter, req *http.Request) {
	origin, _ := url.Parse(urlTo)
	proxy := httputil.NewSingleHostReverseProxy(origin)

	log.Println("Forwarding to:", origin, req.URL.Path)

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}
